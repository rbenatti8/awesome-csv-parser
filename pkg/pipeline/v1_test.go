package pipeline

import (
	"awesome-csv-parser/internal/concurrency"
	"awesome-csv-parser/pkg/schema"
	"awesome-csv-parser/pkg/schema/formatters"
	"awesome-csv-parser/pkg/schema/sanitizers"
	"context"
	"github.com/stretchr/testify/assert"
	"io"
	"strings"
	"testing"
)

type MockFileWriter struct {
	ContentBuilder strings.Builder
}

func (m *MockFileWriter) Write(p []byte) (n int, err error) {
	return m.ContentBuilder.Write(p)
}

func (m *MockFileWriter) Close() error {
	return nil
}

type MockFileOpener struct {
	Reader *strings.Reader
}

func (m *MockFileOpener) Read(p []byte) (n int, err error) {
	return m.Reader.Read(p)
}

func (m *MockFileOpener) Close() error {
	return nil
}

func TestNewV2(t *testing.T) {
	type args struct {
		sourcePath string
		schema     *schema.V1
		opts       []Option
	}
	tests := []struct {
		want *V1
		name string
		args args
	}{
		{
			name: "should return a new V1 instance",
			args: args{
				sourcePath: "sourcePath",
				schema:     &schema.V1{},
				opts:       nil,
			},
			want: &V1{
				sourcePath: "sourcePath",
				config: &Config{
					NumWorkers:        1,
					BatchSize:         1,
					FileOpener:        defaultFileOpener,
					SuccessFileWriter: defaultFileWriter,
					FailureFileWriter: defaultFileWriter,
					SuccessPath:       "../results/success",
					FailurePath:       "../results/failure",
				},
				schema:  &schema.V1{},
				metrics: &Metrics{TotalRowsProcessed: 0, TotalRowsFailed: 0},
			},
		},
		{
			name: "should return a new V1 instance with options",
			args: args{
				sourcePath: "sourcePath",
				schema:     &schema.V1{},
				opts: []Option{
					WithNumWorkers(2),
					WithBatchSize(2),
					WithSuccessPath("successPath"),
					WithFailurePath("failurePath"),
					WithFailureFileWriter(func(name string) (io.WriteCloser, error) {
						return &MockFileWriter{}, nil
					}),
					WithSuccessFileWriter(func(name string) (io.WriteCloser, error) {
						return &MockFileWriter{}, nil
					}),
					WithFileOpener(func(name string) (io.ReadCloser, error) {
						return &MockFileOpener{}, nil
					}),
				},
			},
			want: &V1{
				sourcePath: "sourcePath",
				config: &Config{
					NumWorkers: 2,
					BatchSize:  2,
					FileOpener: func(name string) (io.ReadCloser, error) {
						return &MockFileOpener{}, nil
					},
					SuccessFileWriter: func(name string) (io.WriteCloser, error) {
						return &MockFileWriter{}, nil
					},
					FailureFileWriter: func(name string) (io.WriteCloser, error) {
						return &MockFileWriter{}, nil
					},
					SuccessPath: "successPath",
					FailurePath: "failurePath",
				},
				schema:  &schema.V1{},
				metrics: &Metrics{TotalRowsProcessed: 0, TotalRowsFailed: 0},
			},
		},
	}

	for _, tt := range tests {
		tt = tt

		t.Run(tt.name, func(t *testing.T) {
			if got := NewV1(tt.args.sourcePath, tt.args.schema, tt.args.opts...); got == nil {
				t.Errorf("NewV1() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestV2_Run(t *testing.T) {
	type args struct {
		sourcePath        string
		sourceFileContent string
		schema            *schema.V1
		opts              []Option
	}

	type testCase struct {
		runAssertions func(t *testing.T, successMockFile, failureMockFile *MockFileWriter, metrics *Metrics, err error)
		name          string
		args          args
	}

	testCases := []testCase{
		{
			name: "should write to success file",
			args: args{
				sourcePath: "sourcePath",
				sourceFileContent: `Name,Email,Wage,Number
John Doe,doe@test.com,$10.00,1
Mary Jane,Mary@tes.com,$15,2
Max Topperson,max@test.com,$11,3
Alfred Donald,,$11.5,4
Renato Benatti,renato.benatti.com.br,$11.5,4`,
				schema: &schema.V1{
					ID:   "1",
					Name: "Schema for unit testing",
					TargetFields: []schema.Field{
						{
							Name:         "employee_name",
							SourceFields: []string{"Name"},
							IsRequired:   true,
							IsUnique:     false,
							Sanitizers: []sanitizers.Sanitizer{
								&sanitizers.RemoveSpecialCharacters{},
								&sanitizers.Capitalize{},
							},
						},
						{
							Name:         "employee_salary",
							SourceFields: []string{"Wage"},
							IsRequired:   true,
							IsUnique:     false,
							Sanitizers: []sanitizers.Sanitizer{
								&sanitizers.RemoveSpecialCharacters{
									Exclude: []rune{'.'},
								},
								&sanitizers.RemoveLetters{},
							},
							Formatter: &formatters.Decimal{
								Places: 2,
							},
						},
						{
							Name:         "employee_email",
							SourceFields: []string{"Email"},
							IsRequired:   true,
							IsUnique:     true,
							Sanitizers: []sanitizers.Sanitizer{
								&sanitizers.LowerCase{},
							},
							Formatter: &formatters.Email{},
						},
						{
							Name:         "employee_id",
							SourceFields: []string{"Number"},
							IsRequired:   true,
							IsUnique:     true,
						},
						{
							Name:         "employee_phone",
							SourceFields: []string{"Phone"},
							IsRequired:   false,
							IsUnique:     false,
						},
					},
					ShardedMap: concurrency.NewShardedMap[bool](1024),
				},
				opts: nil,
			},
			runAssertions: func(t *testing.T, successMockFile, failureMockFile *MockFileWriter, metrics *Metrics, err error) {
				assert.NoError(t, err)
				assert.Equal(t, 5, metrics.TotalRowsProcessed)
				assert.Equal(t, 2, metrics.TotalRowsFailed)

				expectedSuccessFileContent := `employee_name,employee_salary,employee_email,employee_id,employee_phone
John Doe,10.00,doe@test.com,1,
Mary Jane,15.00,mary@tes.com,2,
Max Topperson,11.00,max@test.com,3,
`
				expectedErrorFileContent := `Name,Email,Wage,Number,Reason
Alfred Donald,,$11.5,4,field Email is empty
Renato Benatti,renato.benatti.com.br,$11.5,4,renato.benatti.com.br is not a valid email
`

				assert.Equal(t, expectedSuccessFileContent, successMockFile.ContentBuilder.String())
				assert.Equal(t, expectedErrorFileContent, failureMockFile.ContentBuilder.String())
			},
		},
	}

	for _, tt := range testCases {
		tc := tt
		t.Run(tc.name, func(t *testing.T) {
			successMockFile := &MockFileWriter{
				ContentBuilder: strings.Builder{},
			}

			failureMockFile := &MockFileWriter{
				ContentBuilder: strings.Builder{},
			}

			successFileWriter := func(name string) (io.WriteCloser, error) {
				return successMockFile, nil
			}

			failureMockFileWriter := func(name string) (io.WriteCloser, error) {
				return failureMockFile, nil
			}

			sourceFileOpener := func(name string) (io.ReadCloser, error) {
				return &MockFileOpener{
					Reader: strings.NewReader(tc.args.sourceFileContent),
				}, nil
			}

			pipe := NewV1(tc.args.sourcePath, tc.args.schema, WithFileOpener(sourceFileOpener), WithSuccessFileWriter(successFileWriter), WithFailureFileWriter(failureMockFileWriter))
			metrics, err := pipe.Run(context.Background())
			tc.runAssertions(t, successMockFile, failureMockFile, metrics, err)
		})
	}
}
