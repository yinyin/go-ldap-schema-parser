package main

import (
	"log"

	ldapschemaparser "github.com/yinyin/go-ldap-schema-parser"
)

func main() {
	elementStorePath, rootStorePath, outputPath, verbose, err := parseCommandParam()
	if nil != err {
		log.Fatalf("failed on parsing command line parameters: %v", err)
		return
	}
	elementStore := ldapschemaparser.NewLDAPSchemaStore()
	if err = elementStore.ReadFromFile(elementStorePath); nil != err {
		log.Fatalf("ERROR: cannot load element LDAP schema store from [%v]: %v", elementStorePath, err)
		return
	}
	rootStore := ldapschemaparser.NewLDAPSchemaStore()
	if err = rootStore.ReadFromFile(rootStorePath); nil != err {
		log.Fatalf("ERROR: cannot load element LDAP schema store from [%v]: %v", rootStorePath, err)
		return
	}
	if err = rootStore.PullDependentSchema(elementStore, verbose); nil != err {
		log.Fatalf("ERROR: failed on pulling dependent schema: %v", err)
		return
	}
	if err = rootStore.WriteToJSONFile(outputPath); nil != err {
		log.Fatalf("ERROR: cannot write content of LDAP schema store into [%v]: %v", outputPath, err)
	}
}
