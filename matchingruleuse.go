package ldapschemaparser

// MatchingRuleUseSchema represent schema of matching rule uses
type MatchingRuleUseSchema struct {
	NumericOID  string
	Name        []string
	Description string
	Obsolete    bool
	AppliesTo   []string
	Extensions  map[string][]string
}

// NewMatchingRuleUseSchemaViaGenericSchema creates matching rule use schema instance from GenericSchema
func NewMatchingRuleUseSchemaViaGenericSchema(generic *GenericSchema) (result *MatchingRuleUseSchema, err error) {
	if "" == generic.NumericOID {
		err = ErrMissingNumericOID
		return
	}
	appliesTo := generic.getValuesOfParameterizedKeyword("APPLIES")
	if len(appliesTo) == 0 {
		err = &ErrMissingField{
			FieldName: "APPLIES",
		}
		return
	}
	return &MatchingRuleUseSchema{
		NumericOID:  generic.NumericOID,
		Name:        generic.getValuesOfParameterizedKeyword("NAME"),
		Description: generic.getValueOfParameterizedKeyword("DESC"),
		Obsolete:    generic.HasFlagKeyword("OBSOLETE"),
		AppliesTo:   appliesTo,
		Extensions:  generic.fetchExtensionProperties(),
	}, nil
}

// ParseMatchingRuleUseSchema parses object class schema text
func ParseMatchingRuleUseSchema(schemaText string) (matchingRuleUseSchema *MatchingRuleUseSchema, err error) {
	genericSchema, err := Parse(schemaText)
	if nil != err {
		return
	}
	return NewMatchingRuleUseSchemaViaGenericSchema(genericSchema)
}
