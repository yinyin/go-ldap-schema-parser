package main

import (
	"flag"
	"fmt"
)

func parseCommandParam() (rfc4512Path, rfc4517Path, rfc4519Path, rfc4523Path, outputPath string, verbose bool, err error) {
	flag.StringVar(&rfc4512Path, "rfc4512", "", "path to RFC-4512 text file")
	flag.StringVar(&rfc4517Path, "rfc4517", "", "path to RFC-4517 text file")
	flag.StringVar(&rfc4519Path, "rfc4519", "", "path to RFC-4519 text file")
	flag.StringVar(&rfc4519Path, "rfc4523", "", "path to RFC-4523 text file")
	flag.StringVar(&outputPath, "out", "", "path to write into")
	flag.BoolVar(&verbose, "verbose", false, "enable verbose mode")
	flag.Parse()
	if "" == rfc4512Path {
		err = fmt.Errorf("require text file of RFC-4512")
		return
	}
	if "" == rfc4517Path {
		err = fmt.Errorf("require text file of RFC-4517")
		return
	}
	if "" == rfc4519Path {
		err = fmt.Errorf("require text file of RFC-4519")
		return
	}
	if "" == rfc4523Path {
		err = fmt.Errorf("require text file of RFC-4523")
		return
	}
	err = nil
	return
}
