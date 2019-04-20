package ldapschemaparser

import (
	"errors"
	"log"
)

// ErrParseFailed indicate parser stopped at failed state
var ErrParseFailed = errors.New("parsing LDAP schema failed with error parsing state")

// ErrEmptyResult indicate parser resulted an empty result
var ErrEmptyResult = errors.New("parsing LDAP schema failed with empty result")

// SourceRuleType represent source parsing rule
type SourceRuleType int

// constants for source parsing rules
const (
	UnknownRule SourceRuleType = iota
	OIDsRule
)

func undupAppend(s []string, t string) []string {
	for _, v := range s {
		if v == t {
			return s
		}
	}
	s = append(s, t)
	return s
}

// ParameterizedKeyword is keyword with parameters attached
type ParameterizedKeyword struct {
	SourceRule  SourceRuleType
	KeywordText string
	Parameters  []string
}

func newParameterizedKeywordWithParameter(paramText string, sourceRule SourceRuleType) (paramKeyword *ParameterizedKeyword) {
	paramKeyword = &ParameterizedKeyword{
		SourceRule: sourceRule,
	}
	paramKeyword.addParameter(paramText)
	return
}

func (paramKeyword *ParameterizedKeyword) addParameter(paramText string) {
	paramKeyword.Parameters = undupAppend(paramKeyword.Parameters, paramText)
}

func (paramKeyword *ParameterizedKeyword) add(other *ParameterizedKeyword) {
	if (other.SourceRule != paramKeyword.SourceRule) && (other.SourceRule != UnknownRule) && (paramKeyword.SourceRule != UnknownRule) {
		log.Fatalf("cannot add two ParameterizedKeyword with different source rule together: %#v, %#v", paramKeyword, other)
		return
	}
	for _, param := range other.Parameters {
		paramKeyword.addParameter(param)
	}
}

// GenericSchema is generic schema object
type GenericSchema struct {
	NumericOID            string
	FlagKeywords          []string
	ParameterizedKeywords map[string]*ParameterizedKeyword
}

func newGenericSchema() *GenericSchema {
	return &GenericSchema{
		ParameterizedKeywords: make(map[string]*ParameterizedKeyword),
	}
}

func (schema *GenericSchema) addFlagKeywords(keyword string) {
	schema.FlagKeywords = undupAppend(schema.FlagKeywords, keyword)
}

func (schema *GenericSchema) addParameterizedKeyword(keyword string, paramKeyword *ParameterizedKeyword) {
	localParamKeyword, ok := schema.ParameterizedKeywords[keyword]
	if ok {
		localParamKeyword.add(paramKeyword)
	} else {
		paramKeyword.KeywordText = keyword
		schema.ParameterizedKeywords[keyword] = paramKeyword
	}
}

func (schema *GenericSchema) add(other *GenericSchema) {
	for _, kw := range other.FlagKeywords {
		schema.addFlagKeywords(kw)
	}
	for kw, param := range other.ParameterizedKeywords {
		schema.addParameterizedKeyword(kw, param)
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
