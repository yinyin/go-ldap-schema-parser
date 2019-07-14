package ldapschemaparser

import (
	"errors"
	"log"
	"strings"
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
	NumberIDsRule
	OIDWithLengthRule
	QuotedStringsRule
	QuotedStringRule
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

func newParameterizedKeywordWithParameters(paramTexts []string, sourceRule SourceRuleType) (paramKeyword *ParameterizedKeyword) {
	paramKeyword = &ParameterizedKeyword{
		SourceRule: sourceRule,
	}
	for _, paramText := range paramTexts {
		paramKeyword.addParameter(paramText)
	}
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
	if localParamKeyword := schema.ParameterizedKeywords[keyword]; nil != localParamKeyword {
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

func (schema *GenericSchema) getValuesOfParameterizedKeyword(keyword string) []string {
	result := schema.ParameterizedKeywords[keyword]
	if (nil == result) || (0 == len(result.Parameters)) {
		return nil
	}
	return result.Parameters
}

func (schema *GenericSchema) getValueOfParameterizedKeyword(keyword string) string {
	aux := schema.getValuesOfParameterizedKeyword(keyword)
	if nil != aux {
		l := len(aux) - 1
		if l >= 0 {
			return aux[l]
		}
	}
	return ""
}

func (schema *GenericSchema) fetchExtensionProperties() (result map[string][]string) {
	result = make(map[string][]string)
	for keyword, parameters := range schema.ParameterizedKeywords {
		if !isExtensionKeyword(keyword) {
			continue
		}
		if 0 == len(parameters.Parameters) {
			continue
		}
		result[keyword] = parameters.Parameters
	}
	if 0 == len(result) {
		return nil
	}
	return result
}

// HasFlagKeyword checks if given keyword is contained in flag keywords
func (schema *GenericSchema) HasFlagKeyword(keyword string) bool {
	keyword = strings.ToUpper(keyword)
	for _, k := range schema.FlagKeywords {
		if k == keyword {
			return true
		}
	}
	return false
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
