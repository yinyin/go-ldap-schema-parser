package main

import (
	"errors"
	"flag"
)

func parseCommandParam() (elementStorePath, rootStorePath, outputPath string, verbose bool, err error) {
	flag.StringVar(&elementStorePath, "element", "", "path to store for getting schema elements")
	flag.StringVar(&rootStorePath, "root", "", "path to store for getting root elements")
	flag.StringVar(&outputPath, "out", "", "path to write into")
	flag.BoolVar(&verbose, "verbose", false, "enable verbose mode")
	flag.Parse()
	if "" == elementStorePath {
		err = errors.New("require schema element store file (`-element` option)")
		return
	}
	if "" == rootStorePath {
		err = errors.New("require root element store file  (`-root` option)")
		return
	}
	if "" == outputPath {
		err = errors.New("require output file")
		return
	}
	err = nil
	return
}
