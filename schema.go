package ldapschemaparser

import (
	"errors"
)

// ErrParseFailed indicate parser stopped at failed state
var ErrParseFailed = errors.New("parsing LDAP schema failed with error parsing state")

// ErrEmptyResult indicate parser resulted an empty result
var ErrEmptyResult = errors.New("parsing LDAP schema failed with empty result")

// GenericSchema is generic schema object
type GenericSchema struct {
	NumericOID string
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
