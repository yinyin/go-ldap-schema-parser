package ldapschemaparser

// DITContentRuleSchema represent schema of DIT content rule
type DITContentRuleSchema struct {
	NumericOID  string
	Name        []string
	Description string
	Obsolete    bool
	Aux         []string
	Must        []string
	May         []string
	Not         []string
	Extensions  map[string][]string
}

// NewDITContentRuleSchemaViaGenericSchema creates DIT content rule schema instance from GenericSchema
func NewDITContentRuleSchemaViaGenericSchema(generic *GenericSchema) (result *DITContentRuleSchema, err error) {
	if "" == generic.NumericOID {
		err = ErrMissingNumericOID
		return
	}
	return &DITContentRuleSchema{
		NumericOID:  generic.NumericOID,
		Name:        generic.getValuesOfParameterizedKeyword("NAME"),
		Description: generic.getValueOfParameterizedKeyword("DESC"),
		Obsolete:    generic.HasFlagKeyword("OBSOLETE"),
		Aux:         generic.getValuesOfParameterizedKeyword("AUX"),
		Must:        generic.getValuesOfParameterizedKeyword("MUST"),
		May:         generic.getValuesOfParameterizedKeyword("MAY"),
		Not:         generic.getValuesOfParameterizedKeyword("NOT"),
		Extensions:  generic.fetchExtensionProperties(),
	}, nil
}

func (s *DITContentRuleSchema) String() string {
	b := SchemaTextBuilder{}
	b.AppendFragment(s.NumericOID)
	b.AppendQStringSlice("NAME", s.Name)
	b.AppendQString("DESC", s.Description)
	b.AppendFlag("OBSOLETE", s.Obsolete)
	b.AppendOIDSlice("AUX", s.Aux)
	b.AppendOIDSlice("MUST", s.Must)
	b.AppendOIDSlice("MAY", s.May)
	b.AppendOIDSlice("NOT", s.Not)
	b.AppendExtensions(s.Extensions)
	return b.String()
}

// ParseDITContentRuleSchema parses DIT content rule schema text
func ParseDITContentRuleSchema(schemaText string) (ditContentRuleSchema *DITContentRuleSchema, err error) {
	genericSchema, err := Parse(schemaText)
	if nil != err {
		return
	}
	return NewDITContentRuleSchemaViaGenericSchema(genericSchema)
}
