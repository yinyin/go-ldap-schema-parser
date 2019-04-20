package main

import (
	"errors"
	"flag"
	"path/filepath"
)

// ErrInputFileRequired indicates input file path is missing.
var ErrInputFileRequired = errors.New("Input file is required")

// ErrOutputFileRequired indicates output file path is missing.
var ErrOutputFileRequired = errors.New("Output file is required")

func parseCommandParam() (inputFilePath, outputFilePath string, err error) {
	flag.StringVar(&inputFilePath, "in", "", "path to input file")
	flag.StringVar(&outputFilePath, "out", "", "path to output file")
	flag.Parse()
	if "" == inputFilePath {
		err = ErrInputFileRequired
		return
	}
	if inputFilePath, err = filepath.Abs(inputFilePath); nil != err {
		return
	}
	if "" == outputFilePath {
		err = ErrOutputFileRequired
		return
	}
	if outputFilePath, err = filepath.Abs(outputFilePath); nil != err {
		return
	}
	err = nil
	return
}
