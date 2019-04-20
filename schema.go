package ldapschemaparser

import (
	"errors"
)

// ErrParseFailed indicate parser stopped at failed state
var ErrParseFailed = errors.New("parsing LDAP schema failed with error parsing state")

// ErrEmptyResult indicate parser resulted an empty result
var ErrEmptyResult = errors.New("parsing LDAP schema failed with empty result")

func undupAppend(s []string, t string) []string {
	for _, v := range s {
		if v == t {
			return s
		}
	}
	s = append(s, t)
	return s
}

// AttributeWithOIDs is an attribute definition with OIDs as arguments
type AttributeWithOIDs struct {
	KeywordText string
	OIDs        []string
}

func newAttributeWithOIDsWithOID(oidText string) (attr *AttributeWithOIDs) {
	attr = &AttributeWithOIDs{}
	attr.addOID(oidText)
	return
}

func (attr *AttributeWithOIDs) addOID(oidText string) {
	attr.OIDs = undupAppend(attr.OIDs, oidText)
}

func (attr *AttributeWithOIDs) add(other *AttributeWithOIDs) {
	for _, oidText := range other.OIDs {
		attr.addOID(oidText)
	}
}

// GenericSchema is generic schema object
type GenericSchema struct {
	NumericOID   string
	FlagKeywords []string
	OIDKeywords  map[string]*AttributeWithOIDs
}

func newGenericSchema() *GenericSchema {
	return &GenericSchema{
		OIDKeywords: make(map[string]*AttributeWithOIDs),
	}
}

func (schema *GenericSchema) addFlagKeywords(keyword string) {
	schema.FlagKeywords = undupAppend(schema.FlagKeywords, keyword)
}

func (schema *GenericSchema) addAttributeWithOIDs(keyword string, attr *AttributeWithOIDs) {
	localAttr, ok := schema.OIDKeywords[keyword]
	if ok {
		localAttr.add(attr)
	} else {
		attr.KeywordText = keyword
		schema.OIDKeywords[keyword] = attr
	}
}

func (schema *GenericSchema) add(other *GenericSchema) {
	for _, kw := range other.FlagKeywords {
		schema.addFlagKeywords(kw)
	}
	for kw, attr := range other.OIDKeywords {
		schema.addAttributeWithOIDs(kw, attr)
	}
}

// Parse parsing given schema text into generic schema structure
func Parse(schemaText string) (genericSchema *GenericSchema, err error) {
	lexer := newSchemaLexer(schemaText)
	parser := yyNewParser()
	if parser.Parse(lexer) != 0 {
		return nil, ErrParseFailed
	}
	genericSchema = lexer.result // it's a little bit hacky: https://github.com/golang/go/issues/20861
	if nil == genericSchema {
		return nil, ErrEmptyResult
	}
	return genericSchema, nil
}
