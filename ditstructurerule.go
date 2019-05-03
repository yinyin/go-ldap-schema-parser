package ldapschemaparser

// DITStructureRuleSchema represent schema of DIT content rule
type DITStructureRuleSchema struct {
	RuleID      string
	Name        []string
	Description string
	Obsolete    bool
	NameForm    string
	SuperRules  []string
	Extensions  map[string][]string
}

// NewDITStructureRuleSchemaViaGenericSchema creates DIT structure rule schema instance from GenericSchema
func NewDITStructureRuleSchemaViaGenericSchema(generic *GenericSchema) (result *DITStructureRuleSchema, err error) {
	if "" == generic.NumericOID {
		err = ErrMissingRuleID
		return
	}
	nameForm := generic.getValueOfParameterizedKeyword("FORM")
	if "" == nameForm {
		err = &ErrMissingField{
			FieldName: "FORM",
		}
		return
	}
	return &DITStructureRuleSchema{
		RuleID:      generic.NumericOID,
		Name:        generic.getValuesOfParameterizedKeyword("NAME"),
		Description: generic.getValueOfParameterizedKeyword("DESC"),
		Obsolete:    generic.HasFlagKeyword("OBSOLETE"),
		NameForm:    nameForm,
		SuperRules:  generic.getValuesOfParameterizedKeyword("SUP"),
		Extensions:  generic.fetchExtensionProperties(),
	}, nil
}

// ParseDITStructureRuleSchema parses DIT content rule schema text
func ParseDITStructureRuleSchema(schemaText string) (ditStructureRuleSchema *DITStructureRuleSchema, err error) {
	genericSchema, err := Parse(schemaText)
	if nil != err {
		return
	}
	return NewDITStructureRuleSchemaViaGenericSchema(genericSchema)
}
