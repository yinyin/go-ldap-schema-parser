package ldapschemaparser

// ClassKindAbstract, ClassKindStructural and ClassKindAuxiliary indicates
// the kind of object class. (RFC-4512 4.1.1)
const (
	ClassKindAbstract   string = "ABSTRACT"
	ClassKindStructural string = "STRUCTURAL"
	ClassKindAuxiliary  string = "AUXILIARY"
)

// ObjectClassSchema represent schema of object class
type ObjectClassSchema struct {
	NumericOID   string
	Name         []string
	Description  string
	Obsolete     bool
	SuperClasses []string
	ClassKind    string
	Must         []string
	May          []string
	Extensions   map[string][]string
}

// NewObjectClassSchemaViaGenericSchema creates object class schema instance from GenericSchema
func NewObjectClassSchemaViaGenericSchema(generic *GenericSchema) (result *ObjectClassSchema, err error) {
	if "" == generic.NumericOID {
		err = ErrMissingNumericOID
		return
	}
	classKind := ClassKindStructural
	if generic.HasFlagKeyword(ClassKindAbstract) {
		classKind = ClassKindAbstract
	} else if generic.HasFlagKeyword(ClassKindAuxiliary) {
		classKind = ClassKindAuxiliary
	}
	return &ObjectClassSchema{
		NumericOID:   generic.NumericOID,
		Name:         generic.getValuesOfParameterizedKeyword("NAME"),
		Description:  generic.getValueOfParameterizedKeyword("DESC"),
		Obsolete:     generic.HasFlagKeyword("OBSOLETE"),
		SuperClasses: generic.getValuesOfParameterizedKeyword("SUP"),
		ClassKind:    classKind,
		Must:         generic.getValuesOfParameterizedKeyword("MUST"),
		May:          generic.getValuesOfParameterizedKeyword("MAY"),
		Extensions:   generic.fetchExtensionProperties(),
	}, nil
}

func (s *ObjectClassSchema) String() string {
	b := SchemaTextBuilder{}
	b.AppendFragment(s.NumericOID)
	b.AppendQStringSlice("NAME", s.Name)
	b.AppendQString("DESC", s.Description)
	b.AppendFlag("OBSOLETE", s.Obsolete)
	b.AppendOIDSlice("SUP", s.SuperClasses)
	if s.ClassKind != ClassKindStructural {
		b.AppendFragment(s.ClassKind)
	}
	b.AppendOIDSlice("MUST", s.Must)
	b.AppendOIDSlice("MAY", s.May)
	b.AppendExtensions(s.Extensions)
	return b.String()
}

// ParseObjectClassSchema parses object class schema text
func ParseObjectClassSchema(schemaText string) (objectClassSchema *ObjectClassSchema, err error) {
	genericSchema, err := Parse(schemaText)
	if nil != err {
		return
	}
	return NewObjectClassSchemaViaGenericSchema(genericSchema)
}
