package tests

import (
	"awesome-csv-parser/pkg/command"
	"awesome-csv-parser/pkg/employer"
	"bytes"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestIntegration(t *testing.T) {
	t.Parallel()

	schemasFile, err := os.Open("../data/schemas.json")
	if err != nil {
		panic(err)
	}

	employersFile, err := os.Open("../data/employers_schemas.json")
	if err != nil {
		panic(err)
	}

	defer func(sf, ef *os.File) {
		_ = sf.Close()
		_ = ef.Close()
	}(schemasFile, employersFile)

	repo := employer.NewJsonRepository(schemasFile, employersFile)

	c := command.NewRoot(repo)

	got := new(bytes.Buffer)

	c.SetOut(got)
	c.SetErr(got)

	_ = c.Flags().Set("filePath", "../data/roster1.csv")
	_ = c.Flags().Set("employerID", "1")
	_ = c.Flags().Set("successOutputPath", "../results/success")
	_ = c.Flags().Set("failureOutputPath", "../results/failure")

	if err = c.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assert.Contains(t, got.String(), "Total rows processed: 5")
	assert.Contains(t, got.String(), "Total rows failed: 2")
}
