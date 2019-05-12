package ldapschemaparser

import (
	"testing"
)

const sampleLDAPSyntax1 = "( 1.3.6.1.4.1.1466.115.121.1.5 DESC 'Binary' X-NOT-HUMAN-READABLE 'TRUE' )"

func TestLDAPSyntaxSchema_1(t *testing.T) {
	s, err := ParseLDAPSyntaxSchema(sampleLDAPSyntax1)
	if nil != err {
		t.Fatalf("failed on parsing LDAP Syntax sample 1: %v", err)
	}
	v := s.String()
	if v != sampleLDAPSyntax1 {
		t.Errorf("expecting %v but have %v", sampleLDAPSyntax1, v)
	}
}
