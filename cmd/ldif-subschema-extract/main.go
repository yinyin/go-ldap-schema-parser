package main

import (
	"log"
)

func main() {
	ldifPaths, outputPath, verbose, err := parseCommandParam()
	if nil != err {
		log.Fatalf("failed on parsing command line parameters: %v", err)
		return
	}
	for _, ldifPath := range ldifPaths {
		log.Printf("INFO: input LDIF %v", ldifPath)
	}
}
