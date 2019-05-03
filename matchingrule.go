package ldapschemaparser

// MatchingRuleSchema represent schema of matching rule
type MatchingRuleSchema struct {
	NumericOID  string
	Name        []string
	Description string
	Obsolete    bool
	Syntax      string
	Extensions  map[string][]string
}

// NewMatchingRuleSchemaViaGenericSchema creates matching rule schema instance from GenericSchema
func NewMatchingRuleSchemaViaGenericSchema(generic *GenericSchema) (result *MatchingRuleSchema, err error) {
	if "" == generic.NumericOID {
		err = ErrMissingNumericOID
		return
	}
	syntaxNOID := generic.getValueOfParameterizedKeyword("SYNTAX")
	if "" == syntaxNOID {
		err = &ErrMissingField{
			FieldName: "SYNTAX",
		}
		return
	}
	return &MatchingRuleSchema{
		NumericOID:  generic.NumericOID,
		Name:        generic.getValuesOfParameterizedKeyword("NAME"),
		Description: generic.getValueOfParameterizedKeyword("DESC"),
		Obsolete:    generic.HasFlagKeyword("OBSOLETE"),
		Syntax:      syntaxNOID,
		Extensions:  generic.fetchExtensionProperties(),
	}, nil
}

// ParseMatchingRuleSchema parses object class schema text
func ParseMatchingRuleSchema(schemaText string) (matchingRuleSchema *MatchingRuleSchema, err error) {
	genericSchema, err := Parse(schemaText)
	if nil != err {
		return
	}
	return NewMatchingRuleSchemaViaGenericSchema(genericSchema)
}
