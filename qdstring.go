package ldapschemaparser

// QDString escape given string according to `qdstring` rule of RFC-4512 section 4.1
func QDString(v string) string {
	result := make([]rune, 0, len(v)+2)
	result = append(result, '\'')
	for _, ch := range []rune(v) {
		switch ch {
		case '\'':
			result = append(result, '\\')
			result = append(result, '2')
			result = append(result, '7')
		case '\\':
			result = append(result, '\\')
			result = append(result, '5')
			result = append(result, 'C')
		default:
			result = append(result, ch)
		}
	}
	result = append(result, '\'')
	return string(result)
}
