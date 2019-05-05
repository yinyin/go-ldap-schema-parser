package main

import (
	"flag"
	"fmt"
)

func parseCommandParam() (rfc4512Path, rfc4517Path string, err error) {
	flag.StringVar(&rfc4512Path, "rfc4512", "", "path to RFC-4512 text file")
	flag.StringVar(&rfc4517Path, "rfc4517", "", "path to RFC-4517 text file")
	flag.Parse()
	if "" == rfc4512Path {
		err = fmt.Errorf("require text file of RFC-4512")
		return
	}
	if "" == rfc4517Path {
		err = fmt.Errorf("require text file of RFC-4517")
		return
	}
	err = nil
	return
}
