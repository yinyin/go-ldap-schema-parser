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

func (s *MatchingRuleSchema) String() string {
	b := SchemaTextBuilder{}
	b.AppendFragment(s.NumericOID)
	b.AppendQStringSlice("NAME", s.Name)
	b.AppendQString("DESC", s.Description)
	b.AppendFlag("OBSOLETE", s.Obsolete)
	b.AppendBareString("SYNTAX", s.Syntax)
	b.AppendExtensions(s.Extensions)
	return b.String()
}

// ParseMatchingRuleSchema parses matching rule schema text
func ParseMatchingRuleSchema(schemaText string) (matchingRuleSchema *MatchingRuleSchema, err error) {
	genericSchema, err := Parse(schemaText)
	if nil != err {
		return
	}
	return NewMatchingRuleSchemaViaGenericSchema(genericSchema)
}
