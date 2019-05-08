package ldapschemaparser

import (
	"sort"
	"strings"
)

// SchemaTextBuilder supports schema build process
type SchemaTextBuilder struct {
	fragments []string
}

// AppendFragment add given keyword to result
func (b *SchemaTextBuilder) AppendFragment(keyword string) {
	b.fragments = append(b.fragments, keyword)
}

// AppendFlag add flag to result if given value is true
func (b *SchemaTextBuilder) AppendFlag(keyword string, value bool) {
	if !value {
		return
	}
	b.fragments = append(b.fragments, keyword)
}

// AppendQString add quoted string to result
func (b *SchemaTextBuilder) AppendQString(keyword string, value string) {
	if "" == value {
		return
	}
	b.fragments = append(b.fragments, keyword)
	b.fragments = append(b.fragments, QDString(value))
}

// AppendQStringSlice add quoted string slice to result
func (b *SchemaTextBuilder) AppendQStringSlice(keyword string, values []string) {
	l := len(values)
	if 0 == l {
		return
	} else if 1 == l {
		b.AppendQString(keyword, values[0])
		return
	}
	b.fragments = append(b.fragments, keyword)
	b.fragments = append(b.fragments, "(")
	for _, value := range values {
		b.fragments = append(b.fragments, QDString(value))
	}
	b.fragments = append(b.fragments, ")")
}

// AppendBareString add quoted string to result
func (b *SchemaTextBuilder) AppendBareString(keyword string, value string) {
	if "" == value {
		return
	}
	b.fragments = append(b.fragments, keyword)
	b.fragments = append(b.fragments, value)
}

// AppendOIDSlice append OIDs into result
// OIDs are seperated by dollar signs
func (b *SchemaTextBuilder) AppendOIDSlice(keyword string, values []string) {
	l := len(values)
	if 0 == l {
		return
	} else if 1 == l {
		b.AppendBareString(keyword, values[0])
		return
	}
	b.fragments = append(b.fragments, keyword)
	b.fragments = append(b.fragments, "(")
	for idx, value := range values {
		if 0 != idx {
			b.fragments = append(b.fragments, "$")
		}
		b.fragments = append(b.fragments, value)
	}
	b.fragments = append(b.fragments, ")")
}

// AppendExtensions appends extensions to result
func (b *SchemaTextBuilder) AppendExtensions(extensions map[string][]string) {
	if nil == extensions {
		return
	}
	l := len(extensions)
	if 0 == l {
		return
	}
	k := make([]string, 0, l)
	for extKey := range extensions {
		k = append(k, extKey)
	}
	sort.Strings(k)
	for _, extKey := range k {
		extValue := extensions[extKey]
		b.AppendQStringSlice(extKey, extValue)
	}
}

func (b *SchemaTextBuilder) String() string {
	return "( " + strings.Join(b.fragments, " ") + " )"
}
