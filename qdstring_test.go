package ldapschemaparser

import (
	"testing"
)

func TestQDString(t *testing.T) {
	var result string
	result = QDString("abc")
	if result != "'abc'" {
		t.Errorf("expecting 'abc', got: %v", result)
	}
	result = QDString("'quote'")
	if result != "'\\27quote\\27'" {
		t.Errorf("expecting '\\27quote\\27', got: %v", result)
	}
	result = QDString("abc of '12")
	if result != "'abc of \\2712'" {
		t.Errorf("expecting 'abc of \\2712', got: %v", result)
	}
	result = QDString("abc of \\data")
	if result != "'abc of \\5Cdata'" {
		t.Errorf("expecting 'abc of \\5Cdata', got: %v", result)
	}
	result = QDString("abc of '12 \\\\data")
	if result != "'abc of \\2712 \\5C\\5Cdata'" {
		t.Errorf("expecting 'abc of \\2712 \\5C\\5Cdata', got: %v", result)
	}
}
