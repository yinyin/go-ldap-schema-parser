package main

import (
	"io"
	"log"
)

func loadRFC4512(path string) (err error) {
	fp, err := OpenRFCTextReader(path)
	if nil != err {
		return
	}
	defer fp.Close()
	for {
		l, lineType, err := fp.ReadLine()
		if nil != err {
			if err == io.EOF {
				err = nil
			}
			return err
		}
		log.Printf("> %v: %v", lineType, l)
	}
}

func loadRFC4517(path string) (err error) {
	fp, err := OpenRFCTextReader(path)
	if nil != err {
		return
	}
	defer fp.Close()
	for {
		l, lineType, err := fp.ReadLine()
		if nil != err {
			if err == io.EOF {
				err = nil
			}
			return err
		}
		log.Printf("> %v: %v", lineType, l)
	}
}

func main() {
	rfc4512Path, rfc4517Path, err := parseCommandParam()
	if nil != err {
		log.Fatalf("missing parameter: %v", err)
		return
	}
	loadRFC4512(rfc4512Path)
	loadRFC4517(rfc4517Path)
}
