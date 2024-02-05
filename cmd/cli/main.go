package main

import (
	"awesome-csv-parser/pkg/command"
	"awesome-csv-parser/pkg/employer"
	"fmt"
	"os"
)

func main() {
	schemasFile, err := os.Open("data/schemas.json")
	if err != nil {
		panic(err)
	}

	employersFile, err := os.Open("data/employers_schemas.json")
	if err != nil {
		panic(err)
	}

	defer func(sf, ef *os.File) {
		_ = sf.Close()
		_ = ef.Close()
	}(schemasFile, employersFile)

	repo := employer.NewJsonRepository(schemasFile, employersFile)

	c := command.NewRoot(repo)

	if err = c.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
