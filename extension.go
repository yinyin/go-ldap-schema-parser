package ldapschemaparser

func isExtensionKeyword(keywordText string) bool {
	if len(keywordText) < 2 {
		return false
	}
	prefixCh := keywordText[0:2]
	if (prefixCh == "X-") || (prefixCh == "x-") {
		return true
	}
	return false
}
