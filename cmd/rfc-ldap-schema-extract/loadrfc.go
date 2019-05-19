package main

import (
	"io"
	"log"

	ldapschemaparser "github.com/yinyin/go-ldap-schema-parser"
)

const (
	targetSchemaUnknown int = iota
	targetSchemaSkip
	targetSchemaByIdentifier
	targetSchemaObjectClass
	targetSchemaAttributeType
	targetSchemaMatchingRule
	targetSchemaLDAPSyntax
)

func lookupTargetSchemaByIdentifier(oidTargetMap map[string]int, l string) int {
	genericSchema, err := ldapschemaparser.Parse(l)
	if nil != err {
		log.Printf("WARN: failed on parsing schema for getting target by identifier: %v", err)
		return targetSchemaUnknown
	}
	if "" != genericSchema.NumericOID {
		n := genericSchema.NumericOID
		targetSchema, ok := oidTargetMap[n]
		if !ok {
			log.Printf("WARN: failed on mapping identifier to target: %v", n)
		}
		if targetSchema == targetSchemaByIdentifier {
			log.Printf("WARN: looped target schema type: %v", n)
			targetSchema = targetSchemaUnknown
		}
		return targetSchema
	}
	return targetSchemaUnknown
}

func appendByTargetSchema(targetSchema int, oidTargetMap map[string]int, objectClassSchemas, attributeTypeSchemas, matchingRuleSchema, ldapSyntaxSchemas []string, l string) ([]string, []string, []string, []string) {
	switch targetSchema {
	case targetSchemaSkip:
		break
	case targetSchemaByIdentifier:
		targetSchema = lookupTargetSchemaByIdentifier(oidTargetMap, l)
		log.Printf("remapped: %v <- %v", targetSchema, l)
		return appendByTargetSchema(targetSchema, oidTargetMap, objectClassSchemas, attributeTypeSchemas, matchingRuleSchema, ldapSyntaxSchemas, l)
	case targetSchemaObjectClass:
		objectClassSchemas = append(objectClassSchemas, l)
	case targetSchemaAttributeType:
		attributeTypeSchemas = append(attributeTypeSchemas, l)
	case targetSchemaMatchingRule:
		matchingRuleSchema = append(matchingRuleSchema, l)
	case targetSchemaLDAPSyntax:
		ldapSyntaxSchemas = append(ldapSyntaxSchemas, l)
	default:
		log.Printf("WARN: unknown target schema: %v", targetSchema)
	}
	return objectClassSchemas, attributeTypeSchemas, matchingRuleSchema, ldapSyntaxSchemas
}

func loadRFCContent(path string, verbose bool, schemaModeMap map[string]int, oidTargetMapByChapter map[string]map[string]int, objectClassSchemas, attributeTypeSchemas, matchingRuleSchema, ldapSyntaxSchemas []string) ([]string, []string, []string, []string, error) {
	fp, err := OpenRFCTextReader(path)
	if nil != err {
		return nil, nil, nil, nil, err
	}
	defer fp.Close()
	targetSchema := targetSchemaUnknown
	for {
		l, lineType, err := fp.ReadLine()
		if nil != err {
			if err == io.EOF {
				err = nil
			}
			return objectClassSchemas, attributeTypeSchemas, matchingRuleSchema, ldapSyntaxSchemas, err
		}
		if verbose {
			log.Printf("> %v: %v", lineType, l)
		}
		switch lineType {
		case LineTypeChapter:
			if nextTarget, ok := schemaModeMap[fp.CurrentChapter]; ok {
				targetSchema = nextTarget
			}
		case LineTypeSchema:
			oidTargetMap, ok := oidTargetMapByChapter[fp.CurrentChapter]
			if !ok {
				oidTargetMap = oidTargetMapByChapter["*"]
			}
			objectClassSchemas, attributeTypeSchemas, matchingRuleSchema, ldapSyntaxSchemas = appendByTargetSchema(targetSchema, oidTargetMap, objectClassSchemas, attributeTypeSchemas, matchingRuleSchema, ldapSyntaxSchemas, l)
		}
	}
}

func loadRFC4512(path string, verbose bool, objectClassSchemas, attributeTypeSchemas, matchingRuleSchema, ldapSyntaxSchemas []string) ([]string, []string, []string, []string, error) {
	schemaModeMap := map[string]int{
		"2.4.":   targetSchemaObjectClass,
		"2.6.2.": targetSchemaAttributeType,
		"4.2.":   targetSchemaByIdentifier,
		"4.2.1.": targetSchemaAttributeType,
		"4.3.":   targetSchemaObjectClass,
		"4.4.":   targetSchemaAttributeType,
		"7.":     targetSchemaUnknown,
	}
	oidTargetMapByChapter := map[string]map[string]int{
		"*": {
			"2.5.18.10": targetSchemaAttributeType,
			"2.5.20.1":  targetSchemaObjectClass,
		},
	}
	return loadRFCContent(path, verbose, schemaModeMap, oidTargetMapByChapter, objectClassSchemas, attributeTypeSchemas, matchingRuleSchema, ldapSyntaxSchemas)
}

func loadRFC4517(path string, verbose bool, objectClassSchemas, attributeTypeSchemas, matchingRuleSchema, ldapSyntaxSchemas []string) ([]string, []string, []string, []string, error) {
	schemaModeMap := map[string]int{
		"3.3.1.":  targetSchemaByIdentifier,
		"3.3.2.":  targetSchemaLDAPSyntax,
		"3.3.7.":  targetSchemaByIdentifier,
		"3.3.9.":  targetSchemaLDAPSyntax,
		"3.3.19.": targetSchemaByIdentifier,
		"3.3.21.": targetSchemaLDAPSyntax,
		"3.3.22.": targetSchemaByIdentifier,
		"3.3.23.": targetSchemaLDAPSyntax,
		"3.3.24.": targetSchemaByIdentifier,
		"3.3.25.": targetSchemaLDAPSyntax,
		"4.2.":    targetSchemaMatchingRule,
	}
	oidTargetMapByChapter := map[string]map[string]int{
		"*": {},
		"3.3.1.": {
			"2.5.18.1":                     targetSchemaAttributeType,
			"1.3.6.1.4.1.1466.115.121.1.3": targetSchemaLDAPSyntax,
		},
		"3.3.7.": {
			"2.5.6.4":                       targetSchemaSkip,
			"1.3.6.1.4.1.1466.115.121.1.16": targetSchemaLDAPSyntax,
		},
		"3.3.8.": {
			"2":                             targetSchemaSkip,
			"1.3.6.1.4.1.1466.115.121.1.17": targetSchemaLDAPSyntax,
		},
		"3.3.19.": {
			"2.5.13.2":                      targetSchemaMatchingRule,
			"1.3.6.1.4.1.1466.115.121.1.30": targetSchemaLDAPSyntax,
		},
		"3.3.20.": {
			"2.5.13.16":                     targetSchemaSkip,
			"1.3.6.1.4.1.1466.115.121.1.31": targetSchemaLDAPSyntax,
		},
		"3.3.22.": {
			"2.5.15.3":                      targetSchemaSkip,
			"1.3.6.1.4.1.1466.115.121.1.35": targetSchemaLDAPSyntax,
		},
		"3.3.24.": {
			"2.5.6.2":                       targetSchemaSkip,
			"1.3.6.1.4.1.1466.115.121.1.37": targetSchemaLDAPSyntax,
		},
	}
	return loadRFCContent(path, verbose, schemaModeMap, oidTargetMapByChapter, objectClassSchemas, attributeTypeSchemas, matchingRuleSchema, ldapSyntaxSchemas)
}

func loadRFC4519(path string, verbose bool, objectClassSchemas, attributeTypeSchemas, matchingRuleSchema, ldapSyntaxSchemas []string) ([]string, []string, []string, []string, error) {
	schemaModeMap := map[string]int{
		"2.": targetSchemaAttributeType,
		"3.": targetSchemaObjectClass,
		"7.": targetSchemaUnknown,
	}
	oidTargetMapByChapter := map[string]map[string]int{
		"*": {},
	}
	return loadRFCContent(path, verbose, schemaModeMap, oidTargetMapByChapter, objectClassSchemas, attributeTypeSchemas, matchingRuleSchema, ldapSyntaxSchemas)
}

func loadRFC4523(path string, verbose bool, objectClassSchemas, attributeTypeSchemas, matchingRuleSchema, ldapSyntaxSchemas []string) ([]string, []string, []string, []string, error) {
	schemaModeMap := map[string]int{
		"2.": targetSchemaLDAPSyntax,
		"3.": targetSchemaMatchingRule,
		"4.": targetSchemaAttributeType,
		"5.": targetSchemaObjectClass,
		"7.": targetSchemaUnknown,
	}
	oidTargetMapByChapter := map[string]map[string]int{
		"*": {},
	}
	return loadRFCContent(path, verbose, schemaModeMap, oidTargetMapByChapter, objectClassSchemas, attributeTypeSchemas, matchingRuleSchema, ldapSyntaxSchemas)
}
