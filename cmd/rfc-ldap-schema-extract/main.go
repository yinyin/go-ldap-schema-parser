package main

import (
	"log"
)

func main() {
	rfc4512Path, rfc4517Path, rfc4519Path, verbose, err := parseCommandParam()
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
}
