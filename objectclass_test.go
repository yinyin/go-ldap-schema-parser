package ldapschemaparser

import (
	"testing"
)

const sampleObjectClass1 = "( 0.9.2342.19200300.100.4.6 NAME 'document' SUP top STRUCTURAL " +
	"MUST documentIdentifier MAY ( commonName $ description $ seeAlso $ localityNa" +
	"me $ organizationName $ organizationalUnitName $ documentTitle $ documentVers" +
	"ion $ documentAuthor $ documentLocation $ documentPublisher ) )"

func TestObjectClass_1(t *testing.T) {
	s, err := ParseObjectClassSchema(sampleObjectClass1)
	if nil != err {
		t.Fatalf("failed on parsing Object Class sample 1: %v", err)
	}
	v := s.String()
	if v != sampleObjectClass1 {
		t.Errorf("expecting %v but have %v", sampleObjectClass1, v)
	}
	if len(s.Must) != 1 {
		t.Errorf("expect 1 MUST item but having %d (%v)", len(s.Must), s.Must)
	}
	if len(s.May) != 11 {
		t.Errorf("expect 11 MAY items but having %d (%v)", len(s.May), s.May)
	}
}
