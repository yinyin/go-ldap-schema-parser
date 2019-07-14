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

// GenericSchema create an instance of generic schema from this instance of matching-rule-use schema.
func (s *MatchingRuleUseSchema) GenericSchema() (generic *GenericSchema) {
	generic = newGenericSchema()
	generic.NumericOID = s.NumericOID
	if len(s.Name) > 0 {
		nameParams := newParameterizedKeywordWithParameters(s.Name, QuotedStringsRule)
		generic.addParameterizedKeyword("NAME", nameParams)
	}
	if "" != s.Description {
		descriptionParam := newParameterizedKeywordWithParameter(s.Description, QuotedStringRule)
		generic.addParameterizedKeyword("DESC", descriptionParam)
	}
	if s.Obsolete {
		generic.addFlagKeywords("OBSOLETE")
	}
	appliesToParams := newParameterizedKeywordWithParameters(s.AppliesTo, QuotedStringsRule)
	generic.addParameterizedKeyword("APPLIES", appliesToParams)
	return
}

func (s *MatchingRuleUseSchema) String() string {
	b := SchemaTextBuilder{}
	b.AppendFragment(s.NumericOID)
	b.AppendQStringSlice("NAME", s.Name)
	b.AppendQString("DESC", s.Description)
	b.AppendFlag("OBSOLETE", s.Obsolete)
	b.AppendOIDSlice("APPLIES", s.AppliesTo)
	b.AppendExtensions(s.Extensions)
	return b.String()
}

// ParseMatchingRuleUseSchema parses matching rule use schema text
func ParseMatchingRuleUseSchema(schemaText string) (matchingRuleUseSchema *MatchingRuleUseSchema, err error) {
	genericSchema, err := Parse(schemaText)
	if nil != err {
		return
	}
	return NewMatchingRuleUseSchemaViaGenericSchema(genericSchema)
}
