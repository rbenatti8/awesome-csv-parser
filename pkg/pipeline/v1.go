package pipeline

import (
	"awesome-csv-parser/pkg/schema"
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"
	"unicode/utf8"
)

type Row struct {
	HeadersMap schema.HeaderMap
	Values     []string
}

type ProcessedRow struct {
	Err     error
	Headers []string
	Values  []string
}

type OpenFileFunc func(name string) (io.ReadCloser, error)
type CreateFileFunc func(name string) (io.WriteCloser, error)

func defaultFileOpener(name string) (io.ReadCloser, error) {
	return os.Open(name)
}

func defaultFileWriter(name string) (io.WriteCloser, error) {
	return os.Create(name)
}

type Config struct {
	FileOpener        OpenFileFunc
	SuccessFileWriter CreateFileFunc
	FailureFileWriter CreateFileFunc
	SuccessPath       string
	FailurePath       string
	NumWorkers        int
	BatchSize         int
}

type Option func(*Config)

type Metrics struct {
	TotalRowsProcessed int
	TotalRowsFailed    int
	TimeTaken          time.Duration
}

func WithNumWorkers(numWorkers int) Option {
	return func(c *Config) {
		c.NumWorkers = numWorkers
	}
}

func WithBatchSize(batchSize int) Option {
	return func(c *Config) {
		c.BatchSize = batchSize
	}
}

func WithSuccessPath(successPath string) Option {
	return func(c *Config) {
		c.SuccessPath = successPath
	}
}

func WithFailurePath(failurePath string) Option {
	return func(c *Config) {
		c.FailurePath = failurePath
	}
}

func WithFileOpener(fileOpener OpenFileFunc) Option {
	return func(c *Config) {
		c.FileOpener = fileOpener
	}
}

func WithSuccessFileWriter(fileWriter CreateFileFunc) Option {
	return func(c *Config) {
		c.SuccessFileWriter = fileWriter
	}
}

func WithFailureFileWriter(fileWriter CreateFileFunc) Option {
	return func(c *Config) {
		c.FailureFileWriter = fileWriter
	}
}

type V1 struct {
	config     *Config
	schema     *schema.V1
	metrics    *Metrics
	sourcePath string
}

func NewV1(sourcePath string, schema *schema.V1, opts ...Option) *V1 {
	c := &Config{
		NumWorkers:        1,
		BatchSize:         1,
		FileOpener:        defaultFileOpener,
		SuccessFileWriter: defaultFileWriter,
		FailureFileWriter: defaultFileWriter,
		SuccessPath:       "results/success",
		FailurePath:       "results/failure",
	}

	for _, opt := range opts {
		opt(c)
	}

	return &V1{
		schema:     schema,
		sourcePath: sourcePath,
		config:     c,
		metrics: &Metrics{
			TotalRowsProcessed: 0,
			TotalRowsFailed:    0,
		},
	}
}

func (p *V1) Run(ctx context.Context) (*Metrics, error) {
	t := time.Now()

	sourceFile, err := p.openSourceFile()
	if err != nil {
		return nil, err
	}

	successOutputFile, failureOutputFile, err := p.createResultFiles(fmt.Sprintf("%d.csv", t.UnixNano()))
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = sourceFile.Close()
		_ = successOutputFile.Close()
		_ = failureOutputFile.Close()
	}()

	successWriter, failureWriter := csv.NewWriter(successOutputFile), csv.NewWriter(failureOutputFile)
	defer func() {
		successWriter.Flush()
		failureWriter.Flush()
	}()

	rowsCh, headers, err := p.reader(ctx, sourceFile)
	if err != nil {
		return nil, err
	}

	failureHeaders := append(headers, "Reason")
	successHeaders := p.schema.Headers()

	_ = failureWriter.Write(failureHeaders)
	_ = successWriter.Write(successHeaders)

	processedRowsCh := make([]<-chan []ProcessedRow, p.config.NumWorkers)

	for i := 0; i < p.config.NumWorkers; i++ {
		processedRowsCh[i] = p.processor(ctx, rowsCh)
	}

	for processedRows := range p.merger(ctx, processedRowsCh...) {
		p.writeResults(processedRows, failureWriter, successWriter)
	}

	p.metrics.TimeTaken = time.Since(t)

	return p.metrics, nil

}

func (p *V1) writeResults(processedRows []ProcessedRow, failureWriter *csv.Writer, successWriter *csv.Writer) {
	for _, row := range processedRows {
		p.metrics.TotalRowsProcessed++

		if row.Err != nil {
			_ = failureWriter.Write(append(row.Values, row.Err.Error()))
			p.metrics.TotalRowsFailed++

			continue
		}

		_ = successWriter.Write(row.Values)
	}
}

func (p *V1) openSourceFile() (io.ReadCloser, error) {
	return p.config.FileOpener(p.sourcePath)
}

func (p *V1) createResultFiles(name string) (successFile io.WriteCloser, failureFile io.WriteCloser, err error) {
	successPath := fmt.Sprintf("%s/%s", p.config.SuccessPath, name)
	failurePath := fmt.Sprintf("%s/%s", p.config.FailurePath, name)

	if successFile, err = p.config.SuccessFileWriter(successPath); err != nil {
		return nil, nil, err
	}

	if failureFile, err = p.config.FailureFileWriter(failurePath); err != nil {
		return nil, nil, err
	}

	return successFile, failureFile, nil
}

func (p *V1) reader(ctx context.Context, file io.Reader) (<-chan []Row, []string, error) {
	r := csv.NewReader(file)

	out := make(chan []Row)

	headers, err := r.Read()
	if err != nil {
		return nil, nil, err
	}

	hm := buildHeadersMap(headers)
	buffer := make([]Row, 0, p.config.BatchSize)

	go func(buffer *[]Row) {
		defer close(out)
		var row []string

		for {
			select {
			case <-ctx.Done():
				return
			default:
				row, err = r.Read()
				if err == io.EOF {
					out <- *buffer
					return
				}

				if err != nil {
					fmt.Println(fmt.Sprintf("error reading row: %s", err))
					continue
				}

				if len(*buffer) == p.config.BatchSize {
					out <- *buffer

					*buffer = make([]Row, 0, p.config.BatchSize)
				}

				*buffer = append(*buffer, Row{
					HeadersMap: hm,
					Values:     row,
				})
			}
		}
	}(&buffer)

	return out, headers, nil
}

func (p *V1) processor(ctx context.Context, rowBatch <-chan []Row) <-chan []ProcessedRow {
	out := make(chan []ProcessedRow)

	processedRows := make([]ProcessedRow, 0, p.config.BatchSize)

	go func(processedRows *[]ProcessedRow) {
		defer close(out)

		targetHeaders := p.schema.Headers()

		for {
			select {
			case <-ctx.Done():
				return
			case row, open := <-rowBatch:
				if !open {
					return
				}

				for _, r := range row {
					fields, err := p.schema.Build(r.Values, r.HeadersMap)

					processedRow := ProcessedRow{
						Headers: targetHeaders,
						Values:  fields,
					}

					if err != nil {
						processedRow = ProcessedRow{
							Headers: r.HeadersMap.Keys(),
							Values:  r.Values,
							Err:     err,
						}
					}

					*processedRows = append(*processedRows, processedRow)
				}

				out <- *processedRows

				*processedRows = make([]ProcessedRow, 0, p.config.BatchSize)
			}
		}
	}(&processedRows)

	return out
}

func (p *V1) merger(ctx context.Context, inputs ...<-chan []ProcessedRow) <-chan []ProcessedRow {
	out := make(chan []ProcessedRow)

	wg := new(sync.WaitGroup)
	multiplexer := func(p <-chan []ProcessedRow) {
		defer wg.Done()

		for in := range p {
			select {
			case <-ctx.Done():
				return
			case out <- in:
			}
		}
	}

	wg.Add(len(inputs))
	for _, in := range inputs {
		go multiplexer(in)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

func buildHeadersMap(headers []string) schema.HeaderMap {
	headersMap := schema.NewHeaderMap(len(headers))

	for i := 0; i < len(headers); i++ {
		h := removeBOM(headers[i])

		headersMap.Set(h, i)
	}

	return headersMap
}

func removeBOM(input string) string {
	if strings.HasPrefix(input, "\xef\xbb\xbf") {
		_, size := utf8.DecodeRuneInString(input)
		input = input[size:]
	}

	return input
}
