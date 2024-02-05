@PHONY: build
build:
	go build -o bin/awesome-csv-parser ./cmd/cli/

@PHONY: run
run: build
	./bin/awesome-csv-parser

@PHONY: test
test:
	go test -v ./... -race -parallel 4

@PHONY: test-cover
test-cover:
	go test -coverprofile=coverage.out ./... -race -parallel 4
	go tool cover -html=coverage.out

@PHONY: clean-results
clean-results:
	rm -rf results/success/*
	rm -rf results/failure/*

@PHONY: roster1
roster1:
	./bin/awesome-csv-parser --employerID 1 --filePath data/roster1.csv

@PHONY: roster2
roster2:
	./bin/awesome-csv-parser --employerID 2 --filePath data/roster2.csv

@PHONY: roster3
roster3:
	./bin/awesome-csv-parser --employerID 3 --filePath data/roster3.csv

@PHONY: roster4
roster4:
	./bin/awesome-csv-parser --employerID 4 --filePath data/roster4.csv
