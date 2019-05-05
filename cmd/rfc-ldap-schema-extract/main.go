package main

import (
	"io"
	"log"

	ldapschemaparser "github.com/yinyin/go-ldap-schema-parser"
)

const (
	targetSchemaUnknown int = iota
	targetSchemaByName
	targetSchemaObjectClass
	targetSchemaAttributeType
)

func lookupTargetSchemaByName(nameTargetMap map[string]int, l string) int {
	genericSchema, err := ldapschemaparser.Parse(l)
	if nil != err {
		log.Printf("WARN: failed on parsing schema for getting target by name: %v", err)
		return targetSchemaUnknown
	}
	if aux := genericSchema.ParameterizedKeywords["NAME"]; (nil != aux) && (0 < len(aux.Parameters)) {
		n := aux.Parameters[0]
		return nameTargetMap[n]
	}
	return targetSchemaUnknown
}

func appendByTargetSchema(targetSchema int, nameTargetMap map[string]int, objectClassSchemas, attributeTypeSchemas []string, l string) ([]string, []string) {
	switch targetSchema {
	case targetSchemaByName:
		targetSchema = lookupTargetSchemaByName(nameTargetMap, l)
		return appendByTargetSchema(targetSchema, nameTargetMap, objectClassSchemas, attributeTypeSchemas, l)
	case targetSchemaObjectClass:
		objectClassSchemas = append(objectClassSchemas, l)
	case targetSchemaAttributeType:
		attributeTypeSchemas = append(attributeTypeSchemas, l)
	default:
		log.Printf("WARN: unknown target schema: %v", targetSchema)
	}
	return objectClassSchemas, attributeTypeSchemas
}

func loadRFC4512(path string, objectClassSchemas, attributeTypeSchemas []string) ([]string, []string, error) {
	fp, err := OpenRFCTextReader(path)
	if nil != err {
		return nil, nil, err
	}
	defer fp.Close()
	schemaModeMap := map[string]int{
		"2.4.":   targetSchemaObjectClass,
		"2.6.2.": targetSchemaAttributeType,
		"4.2.":   targetSchemaByName,
		"4.2.1.": targetSchemaAttributeType,
		"4.3.":   targetSchemaObjectClass,
		"4.4.":   targetSchemaAttributeType,
		"7.":     targetSchemaUnknown,
	}
	nameTargetMap := map[string]int{
		"subschemaSubentry": targetSchemaAttributeType,
		"subschema":         targetSchemaObjectClass,
	}
	targetSchema := targetSchemaUnknown
	for {
		l, lineType, err := fp.ReadLine()
		if nil != err {
			if err == io.EOF {
				err = nil
			}
			return objectClassSchemas, attributeTypeSchemas, err
		}
		log.Printf("> %v: %v", lineType, l)
		switch lineType {
		case LineTypeChapter:
			if nextTarget, ok := schemaModeMap[fp.CurrentChapter]; ok {
				targetSchema = nextTarget
			}
		case LineTypeSchema:
			objectClassSchemas, attributeTypeSchemas = appendByTargetSchema(targetSchema, nameTargetMap, objectClassSchemas, attributeTypeSchemas, l)
		}
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
	objectClassSchemas, attributeTypeSchemas, err := loadRFC4512(rfc4512Path, nil, nil)
	log.Printf("INFO: load RFC 4512: err=%v", err)
	log.Printf("** Object Class:")
	for _, l := range objectClassSchemas {
		log.Print(l)
	}
	log.Printf("** Attribute Type:")
	for _, l := range attributeTypeSchemas {
		log.Print(l)
	}
	loadRFC4517(rfc4517Path)
}
