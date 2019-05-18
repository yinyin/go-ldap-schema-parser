package main

import (
	"log"
	"os"

	ldapschemaparser "github.com/yinyin/go-ldap-schema-parser"
)

func writeToFile(outputPath string, objectClassSchemas, attributeTypeSchemas, matchingRuleSchema, ldapSyntaxSchemas []string) (err error) {
	store := ldapschemaparser.NewLDAPSchemaStore()
	if err = store.ReadFromFile(outputPath); !os.IsNotExist(err) {
		return
	}
	for _, l := range ldapSyntaxSchemas {
		if err = store.AddLDAPSyntaxSchemaText(l); nil != err {
			log.Printf("ERR: failed on importing LDAP syntax schema text to store: %v", l)
			return
		}
	}
	for _, l := range matchingRuleSchema {
		if err = store.AddMatchingRuleSchemaText(l); nil != err {
			log.Printf("ERR: failed on importing matching rule schema text to store: %v", l)
			return
		}
	}
	for _, l := range attributeTypeSchemas {
		if err = store.AddAttributeTypeSchemaText(l); nil != err {
			log.Printf("ERR: failed on importing attribute type schema text to store: %v", l)
			return
		}
	}
	for _, l := range objectClassSchemas {
		if err = store.AddObjectClassSchemaText(l); nil != err {
			log.Printf("ERR: failed on importing object class schema text to store: %v", l)
			return
		}
	}
	return store.WriteToFile(outputPath)
}

func main() {
	rfc4512Path, rfc4517Path, rfc4519Path, outputPath, verbose, err := parseCommandParam()
	if nil != err {
		log.Fatalf("missing parameter: %v", err)
		return
	}
	objectClassSchemas, attributeTypeSchemas, matchingRuleSchema, ldapSyntaxSchemas, err := loadRFC4512(rfc4512Path, verbose, nil, nil, nil, nil)
	log.Printf("INFO: load RFC 4512: err=%v", err)
	objectClassSchemas, attributeTypeSchemas, matchingRuleSchema, ldapSyntaxSchemas, err = loadRFC4517(rfc4517Path, verbose, objectClassSchemas, attributeTypeSchemas, matchingRuleSchema, ldapSyntaxSchemas)
	log.Printf("INFO: load RFC 4517: err=%v", err)
	objectClassSchemas, attributeTypeSchemas, matchingRuleSchema, ldapSyntaxSchemas, err = loadRFC4519(rfc4519Path, verbose, objectClassSchemas, attributeTypeSchemas, matchingRuleSchema, ldapSyntaxSchemas)
	log.Printf("INFO: load RFC 4519: err=%v", err)
	log.Printf("** Object Class (%d):", len(objectClassSchemas))
	for _, l := range objectClassSchemas {
		log.Print(l)
	}
	log.Printf("** Attribute Type (%d):", len(attributeTypeSchemas))
	for _, l := range attributeTypeSchemas {
		log.Print(l)
	}
	log.Printf("** Matching Rule (%d):", len(matchingRuleSchema))
	for _, l := range matchingRuleSchema {
		log.Print(l)
	}
	log.Printf("** LDAP Syntax (%d):", len(ldapSyntaxSchemas))
	for _, l := range ldapSyntaxSchemas {
		log.Print(l)
	}
	if "" != outputPath {
		err = writeToFile(outputPath, objectClassSchemas, attributeTypeSchemas, matchingRuleSchema, ldapSyntaxSchemas)
		if nil != err {
			log.Fatalf("write to file failed: %v", err)
		}
	}
}
