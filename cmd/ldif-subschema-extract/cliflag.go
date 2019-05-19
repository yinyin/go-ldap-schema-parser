package main

import (
	"errors"
	"flag"
)

func parseCommandParam() (ldifPaths []string, outputPath string, verbose bool, err error) {
	flag.StringVar(&outputPath, "out", "", "path to write into")
	flag.BoolVar(&verbose, "verbose", false, "enable verbose mode")
	flag.Parse()
	ldifPaths = flag.Args()
	if 0 == len(ldifPaths) {
		err = errors.New("require input LDIF files")
		return
	}
	err = nil
	return
}
