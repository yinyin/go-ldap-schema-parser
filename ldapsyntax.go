package ldapschemaparser

// LDAPSyntaxSchema represent schema of LDAP syntax
type LDAPSyntaxSchema struct {
	NumericOID  string
	Description string
	Extensions  map[string][]string
}

// NewLDAPSyntaxSchemaViaGenericSchema creates matching rule use schema instance from GenericSchema
func NewLDAPSyntaxSchemaViaGenericSchema(generic *GenericSchema) (result *LDAPSyntaxSchema, err error) {
	if "" == generic.NumericOID {
		err = ErrMissingNumericOID
		return
	}
	return &LDAPSyntaxSchema{
		NumericOID:  generic.NumericOID,
		Description: generic.getValueOfParameterizedKeyword("DESC"),
		Extensions:  generic.fetchExtensionProperties(),
	}, nil
}

// ParseLDAPSyntaxSchema parses LDAP syntax schema text
func ParseLDAPSyntaxSchema(schemaText string) (ldapSyntaxSchema *LDAPSyntaxSchema, err error) {
	genericSchema, err := Parse(schemaText)
	if nil != err {
		return
	}
	return NewLDAPSyntaxSchemaViaGenericSchema(genericSchema)
}
