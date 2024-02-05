package command

import (
	"awesome-csv-parser/pkg/employer"
	"awesome-csv-parser/pkg/pipeline"
	"context"
	"fmt"
	"github.com/spf13/cobra"
)

type employerRepository interface {
	GetByID(id string) (employer *employer.Employer, err error)
}

type Root struct {
	*cobra.Command
	repo employerRepository
}

// NewRoot creates a new root command
func NewRoot(repo employerRepository) *Root {
	cmd := &Root{
		Command: &cobra.Command{
			Use:   "awesome-csv-parser",
			Short: "A CSV parser that can handle large files",
		},
		repo: repo,
	}

	cmd.Flags().String("filePath", "", "File path")
	cmd.Flags().String("employerID", "", "Employer ID")
	cmd.Flags().String("successOutputPath", "results/success", "Success output path")
	cmd.Flags().String("failureOutputPath", "results/failure", "Failure output path")
	cmd.Flags().Int("batchSize", 50_000, "Batch size")
	cmd.Flags().Int("numWorkers", 10, "Number of workers")

	_ = cmd.MarkFlagRequired("filePath")
	_ = cmd.MarkFlagRequired("employerID")

	return cmd
}

// Execute runs the root command
func (r *Root) Execute() error {
	r.Run = r.run
	return r.Command.Execute()
}

func (r *Root) run(cmd *cobra.Command, _ []string) {
	filePath, err := cmd.Flags().GetString("filePath")
	if err != nil {
		cmd.PrintErrln(err)
		return
	}

	employerID, err := cmd.Flags().GetString("employerID")
	if err != nil {
		cmd.PrintErrln(err)
		return
	}

	em, err := r.repo.GetByID(employerID)
	if err != nil {
		cmd.PrintErrln(err)
		return
	}

	opts := buildPipeOpts(cmd)
	pipe := pipeline.NewV1(filePath, em.Schema, opts...)

	metrics, err := pipe.Run(context.Background())
	if err != nil {
		cmd.Println("Error running pipeline: ", err)
		cmd.PrintErrln(err)
		return
	}

	cmd.Println(fmt.Sprintf("Total rows processed: %d", metrics.TotalRowsProcessed))
	cmd.Println(fmt.Sprintf("Total rows failed: %d", metrics.TotalRowsFailed))
	cmd.Println(fmt.Sprintf("Total time taken: %s", metrics.TimeTaken))
}

func buildPipeOpts(cmd *cobra.Command) []pipeline.Option {
	opts := make([]pipeline.Option, 0, 4)

	if numWorkers, err := cmd.Flags().GetInt("numWorkers"); err == nil && numWorkers > 0 {
		opts = append(opts, pipeline.WithNumWorkers(numWorkers))
	}

	if batchSize, err := cmd.Flags().GetInt("batchSize"); err == nil && batchSize > 0 {
		opts = append(opts, pipeline.WithBatchSize(batchSize))
	}

	if successOutputPath, err := cmd.Flags().GetString("successOutputPath"); err == nil && successOutputPath != "" {
		opts = append(opts, pipeline.WithSuccessPath(successOutputPath))
	}

	if failureOutputPath, err := cmd.Flags().GetString("failureOutputPath"); err == nil && failureOutputPath != "" {
		opts = append(opts, pipeline.WithFailurePath(failureOutputPath))
	}

	return opts
}
