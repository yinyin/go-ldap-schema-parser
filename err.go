package ldapschemaparser

import (
	"errors"
	"fmt"
)

// ErrMissingNumericOID indicates numeric OID is required but not given.
var ErrMissingNumericOID = errors.New("NumericOID is required")

// ErrMissingField represents a required field is missing.
type ErrMissingField struct {
	FieldName string
}

func (missingField *ErrMissingField) Error() string {
	return fmt.Sprintf("%s is required field", missingField.FieldName)
}
