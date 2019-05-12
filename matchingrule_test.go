package ldapschemaparser

import (
	"testing"
)

const sampleMatchRule1 = "( 2.5.13.34 NAME 'certificateExactMatch' SYNTAX 1.3.6.1.1.15.1 )"

func TestMatchRuleSchema_1(t *testing.T) {
	s, err := ParseMatchingRuleSchema(sampleMatchRule1)
	if nil != err {
		t.Fatalf("failed on parsing Match Rule sample 1: %v", err)
	}
	v := s.String()
	if v != sampleMatchRule1 {
		t.Errorf("expecting %v but have %v", sampleMatchRule1, v)
	}
}
