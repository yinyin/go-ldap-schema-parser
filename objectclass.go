package ldapschemaparser

// ClassKindAbstract, ClassKindStructural and ClassKindAuxiliary indicates
// the kind of object class. (RFC-4512 4.1.1)
const (
	ClassKindAbstract   = "ABSTRACT"
	ClassKindStructural = "STRUCTURAL"
	ClassKindAuxiliary  = "AUXILIARY"
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
