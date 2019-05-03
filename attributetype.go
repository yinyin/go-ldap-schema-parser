package ldapschemaparser

// AttributeUsageUserApplications, AttributeUsageDirectoryOperation,
// AttributeUsageDistributedOperation and AttributeUsageDSAOperation
// are usage of attribute types. (RFC-4512 4.1.2)
const (
	AttributeUsageUserApplications     string = "userApplications"
	AttributeUsageDirectoryOperation   string = "directoryOperation"
	AttributeUsageDistributedOperation string = "distributedOperation"
	AttributeUsageDSAOperation         string = "dSAOperation"
)

// AttributeTypeSchema represent schema of attribute type
type AttributeTypeSchema struct {
	NumericOID         string
	Name               []string
	Description        string
	Obsolete           bool
	SuperType          string
	Equality           string
	Ordering           string
	SubString          string
	Syntax             string
	SyntaxOID          string
	SyntaxLength       int32
	SingleValue        bool
	Collective         bool
	NoUserModification bool
	Usage              string
	Extensions         map[string][]string
}

// NewAttributeTypeSchemaViaGenericSchema creates attribute type schema instance from GenericSchema
func NewAttributeTypeSchemaViaGenericSchema(generic *GenericSchema) (result *AttributeTypeSchema, err error) {
	if "" == generic.NumericOID {
		err = ErrMissingNumericOID
		return
	}
	syntaxNOIDLen := generic.getValueOfParameterizedKeyword("SYNTAX")
	syntaxNOID, syntaxLen := parseOIDLength(syntaxNOIDLen)
	attrUsage := generic.getValueOfParameterizedKeyword("USAGE")
	if attrUsage == "" {
		attrUsage = AttributeUsageUserApplications
	}
	return &AttributeTypeSchema{
		NumericOID:         generic.NumericOID,
		Name:               generic.getValuesOfParameterizedKeyword("NAME"),
		Description:        generic.getValueOfParameterizedKeyword("DESC"),
		Obsolete:           generic.HasFlagKeyword("OBSOLETE"),
		SuperType:          generic.getValueOfParameterizedKeyword("SUP"),
		Equality:           generic.getValueOfParameterizedKeyword("EQUALITY"),
		Ordering:           generic.getValueOfParameterizedKeyword("ORDERING"),
		SubString:          generic.getValueOfParameterizedKeyword("SUBSTR"),
		Syntax:             syntaxNOIDLen,
		SyntaxOID:          syntaxNOID,
		SyntaxLength:       syntaxLen,
		SingleValue:        generic.HasFlagKeyword("SINGLE-VALUE"),
		Collective:         generic.HasFlagKeyword("COLLECTIVE"),
		NoUserModification: generic.HasFlagKeyword("NO-USER-MODIFICATION"),
		Usage:              attrUsage,
		Extensions:         generic.fetchExtensionProperties(),
	}, nil
}

// ParseAttributeTypeSchema parses attribute type schema text
func ParseAttributeTypeSchema(schemaText string) (attributeTypeSchema *AttributeTypeSchema, err error) {
	genericSchema, err := Parse(schemaText)
	if nil != err {
		return
	}
	return NewAttributeTypeSchemaViaGenericSchema(genericSchema)
}
