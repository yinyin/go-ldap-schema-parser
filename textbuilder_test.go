package ldapschemaparser

import (
	"testing"
)

func TestSchemaTextBuilder_1(t *testing.T) {
	b := SchemaTextBuilder{}
	b.AppendFragment("1.2.3")
	b.AppendQString("NAME", "test01")
	b.AppendFlag("FLAG01A", true)
	b.AppendFlag("FLAG01B", false)
	b.AppendBareString("SUPA", "sup01")
	b.AppendBareString("SUPB", "")
	b.AppendOIDSlice("OID01A", []string{
		"businessCategory", "description", "destinationIndicator",
	})
	b.AppendOIDSlice("OID01B", []string{
		"postOfficeBox",
	})
	b.AppendQString("QSTR01A", "")
	b.AppendQString("QSTR01B", "beta 21' car")
	b.AppendQStringSlice("QSTR01C", []string{
		"homePhone", "homeTelephoneNumber",
	})
	b.AppendQStringSlice("QSTR01D", nil)
	b.AppendExtensions(map[string][]string{
		"X-TE01A": {"true"},
		"X-TE01B": {"alpha", "beta"},
	})
	result := b.String()
	if result != "( 1.2.3 NAME 'test01' FLAG01A SUPA sup01 "+
		"OID01A ( businessCategory $ description $ destinationIndicator ) "+
		"OID01B postOfficeBox QSTR01B 'beta 21\\27 car' QSTR01C ( 'homePhone' 'homeTelephoneNumber' ) "+
		"X-TE01A 'true' X-TE01B ( 'alpha' 'beta' ) )" {
		t.Errorf("unexpect result: %v", result)
	}
}
