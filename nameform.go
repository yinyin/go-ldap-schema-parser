package ldapschemaparser

// NameFormSchema represent schema of name form
type NameFormSchema struct {
	NumericOID  string
	Name        []string
	Description string
	Obsolete    bool
	ObjectClass string
	Must        []string
	May         []string
	Extensions  map[string][]string
}

// NewNameFormSchemaViaGenericSchema creates name form schema instance from GenericSchema
func NewNameFormSchemaViaGenericSchema(generic *GenericSchema) (result *NameFormSchema, err error) {
	if "" == generic.NumericOID {
		err = ErrMissingNumericOID
		return
	}
	objectClass := generic.getValueOfParameterizedKeyword("OC")
	if "" == objectClass {
		err = &ErrMissingField{
			FieldName: "OC",
		}
		return
	}
	attrMust := generic.getValuesOfParameterizedKeyword("MUST")
	if 0 == len(attrMust) {
		err = &ErrMissingField{
			FieldName: "MUST",
		}
		return
	}
	return &NameFormSchema{
		NumericOID:  generic.NumericOID,
		Name:        generic.getValuesOfParameterizedKeyword("NAME"),
		Description: generic.getValueOfParameterizedKeyword("DESC"),
		Obsolete:    generic.HasFlagKeyword("OBSOLETE"),
		ObjectClass: objectClass,
		Must:        attrMust,
		May:         generic.getValuesOfParameterizedKeyword("MAY"),
		Extensions:  generic.fetchExtensionProperties(),
	}, nil
}

func (s *NameFormSchema) String() string {
	b := SchemaTextBuilder{}
	b.AppendFragment(s.NumericOID)
	b.AppendQStringSlice("NAME", s.Name)
	b.AppendQString("DESC", s.Description)
	b.AppendFlag("OBSOLETE", s.Obsolete)
	b.AppendBareString("OC", s.ObjectClass)
	b.AppendOIDSlice("MUST", s.Must)
	b.AppendOIDSlice("MAY", s.May)
	b.AppendExtensions(s.Extensions)
	return b.String()
}

// ParseNameFormSchema parses name form schema text
func ParseNameFormSchema(schemaText string) (nameFormSchema *NameFormSchema, err error) {
	genericSchema, err := Parse(schemaText)
	if nil != err {
		return
	}
	return NewNameFormSchemaViaGenericSchema(genericSchema)
}
