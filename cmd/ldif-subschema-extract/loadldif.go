package main

import (
	"log"
	"os"

	"github.com/go-ldap/ldif"
	ldap "gopkg.in/ldap.v2"

	ldapschemaparser "github.com/yinyin/go-ldap-schema-parser"
)

func addEntryAttributeToStore(store *ldapschemaparser.LDAPSchemaStore, attr *ldap.EntryAttribute, verbose bool) (err error) {
	switch attr.Name {
	case "ldapSyntaxes":
		for _, schemaText := range attr.Values {
			if verbose {
				log.Printf("ldapSyntaxes: %v", schemaText)
			}
			if nil != store {
				if err = store.AddLDAPSyntaxSchemaText(schemaText); nil != err {
					log.Printf("ERROR: cannot add LDAP syntax schema to store: %v - %v", schemaText, err)
					return err
				}
			}
		}
	case "matchingRules":
		for _, schemaText := range attr.Values {
			if verbose {
				log.Printf("matchingRules: %v", schemaText)
			}
			if nil != store {
				if err = store.AddMatchingRuleSchemaText(schemaText); nil != err {
					log.Printf("ERROR: cannot add matching rule schema to store: %v - %v", schemaText, err)
					return err
				}
			}
		}
	case "matchingRuleUse":
	case "olcAttributeTypes":
		fallthrough
	case "attributeTypes":
		for _, schemaText := range attr.Values {
			if verbose {
				log.Printf("attributeTypes: %v", schemaText)
			}
			if nil != store {
				if err = store.AddAttributeTypeSchemaText(schemaText); nil != err {
					log.Printf("ERROR: cannot add attribute type schema to store: %v - %v", schemaText, err)
					return err
				}
			}
		}
	case "olcObjectClasses":
		fallthrough
	case "objectClasses":
		for _, schemaText := range attr.Values {
			if verbose {
				log.Printf("objectClasses: %v", schemaText)
			}
			if nil != store {
				if err = store.AddObjectClassSchemaText(schemaText); nil != err {
					log.Printf("ERROR: cannot add object class schema to store: %v - %v", schemaText, err)
					return err
				}
			}
		}
	}
	return nil
}

func loadLDIF(store *ldapschemaparser.LDAPSchemaStore, ldifPath string, verbose bool) (err error) {
	fp, err := os.Open(ldifPath)
	if nil != err {
		return
	}
	defer fp.Close()
	var ldifContent ldif.LDIF
	if err = ldif.Unmarshal(fp, &ldifContent); nil != err {
		log.Printf("WARN: failed on unmarshal LDIF %v: %v", ldifPath, err)
	}
	for _, entry := range ldifContent.Entries {
		if nil == entry.Entry {
			continue
		}
		for _, attr := range entry.Entry.Attributes {
			if err = addEntryAttributeToStore(store, attr, verbose); nil != err {
				return err
			}
		}
	}
	return nil
}
