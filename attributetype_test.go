package ldapschemaparser

import (
	"testing"
)

const sampleAttributeType1 = "( 1.3.6.1.1.16.4 NAME 'entryUUID' DESC 'UUID of the entry' EQU" +
	"ALITY UUIDMatch ORDERING UUIDOrderingMatch SYNTAX 1.3.6.1.1.16.1 SINGLE-VALUE " +
	"NO-USER-MODIFICATION USAGE directoryOperation )"

func TestAttributeType_1(t *testing.T) {
	s, err := ParseAttributeTypeSchema(sampleAttributeType1)
	if nil != err {
		t.Fatalf("failed on parsing Attribute Type sample 1: %v", err)
	}
	v := s.String()
	if v != sampleAttributeType1 {
		t.Errorf("expecting %v but have %v", sampleAttributeType1, v)
	}
}
