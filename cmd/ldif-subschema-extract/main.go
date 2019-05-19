package main

import (
	"io"
	"log"

	ldapschemaparser "github.com/yinyin/go-ldap-schema-parser"
)

func main() {
	ldifPaths, outputPath, verbose, err := parseCommandParam()
	if nil != err {
		log.Fatalf("failed on parsing command line parameters: %v", err)
		return
	}
	var store *ldapschemaparser.LDAPSchemaStore
	if "" != outputPath {
		store = ldapschemaparser.NewLDAPSchemaStore()
		if err = store.ReadFromFile(outputPath); nil != err {
			if io.EOF != err {
				log.Fatalf("ERROR: cannot load LDAP schema store from [%v]: %v", outputPath, err)
				return
			}
		}
	}
	for _, ldifPath := range ldifPaths {
		log.Printf("INFO: input LDIF %v", ldifPath)
		if err = loadLDIF(store, ldifPath, verbose); nil != err {
			log.Fatalf("failed on loading LDIF from %v: %v", ldifPath, err)
			return
		}
	}
	log.Printf("INFO: output to: %v", outputPath)
	if "" != outputPath {
		if err = store.WriteToFile(outputPath); nil != err {
			log.Fatalf("ERROR: cannot write content of LDAP schema store into [%v]: %v", outputPath, err)
		}
	}
}
