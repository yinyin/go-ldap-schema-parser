package ldapschemaparser

//go:generate goyacc -o parser.go parser.y

import (
	"log"
	"strings"
	"unicode"
)

const dataEOF = 0

var keywordTypeLookupMap = map[string]int{
	"SUP":  OIDS_ATTR_KEYWORD,
	"MUST": OIDS_ATTR_KEYWORD,
	"MAY":  OIDS_ATTR_KEYWORD,
}

func lookupKeywordType(keywordText string) (string, int) {
	u := strings.ToUpper(keywordText)
	keywordIdentifier, ok := keywordTypeLookupMap[u]
	if ok {
		return u, keywordIdentifier
	}
	return keywordText, KEYWORD
}

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

type schemaLexer struct {
	dataContent  []rune
	dataLength   int
	currentIndex int

	result *GenericSchema
}

func newSchemaLexer(schemaText string) *schemaLexer {
	d := []rune(schemaText)
	return &schemaLexer{
		dataContent: d,
		dataLength:  len(d),
	}
}

func (lexer *schemaLexer) Lex(lval *yySymType) (lexIdentifier int) {
	var result []rune
	startIndex := lexer.currentIndex
	for {
		ch := lexer.next()
		log.Printf("DEBUG: %v", ch)
		if ch == dataEOF {
			break
		}
		switch lexIdentifier {
		case 0:
			if (ch == '(') || (ch == ')') || (ch == '{') || (ch == '}') || (ch == '$') {
				return int(ch)
			}
			if unicode.IsDigit(ch) {
				lexIdentifier = NUMBER
			} else if unicode.IsSpace(ch) {
				lexIdentifier = SPACES
			} else if unicode.IsLetter(ch) {
				lexIdentifier = KEYWORD
			} else if ch == '\'' {
				lexIdentifier = SQSTRING
			} else if ch == '"' {
				lexIdentifier = DQSTRING
			}
		case NUMBER:
			if ch == '.' {
				lexIdentifier = NUMERIC_OID
			} else if !unicode.IsDigit(ch) {
				lexer.putBack()
				lexer.fetchText(lval, startIndex)
				return
			}
		case NUMERIC_OID:
			if (ch != '.') && (!unicode.IsDigit(ch)) {
				lexer.putBack()
				lexer.fetchText(lval, startIndex)
				return
			}
		case SPACES:
			if !unicode.IsSpace(ch) {
				lexer.putBack()
				return
			}
		case KEYWORD:
			if !unicode.IsLetter(ch) && !unicode.IsDigit(ch) && (ch != '-') {
				lexer.putBack()
				// TODO: check if special keyword (eg. NAME, AUX, SUP ...)
				w := lexer.fetchText(lval, startIndex)
				if isExtensionKeyword(w) {
					lexIdentifier = X_KEYWORD
				} else {
					lval.text, lexIdentifier = lookupKeywordType(w)
				}
				return
			}
		case SQSTRING:
			result = lexer.stateTransitQuotedString(lval, result, '\'', ch)
		case DQSTRING:
			result = lexer.stateTransitQuotedString(lval, result, '"', ch)
		}
	}
	return 0
}

func (lexer *schemaLexer) next() rune {
	if lexer.currentIndex >= lexer.dataLength {
		return dataEOF
	}
	ch := lexer.dataContent[lexer.currentIndex]
	lexer.currentIndex++
	return ch
}

func (lexer *schemaLexer) peekString(len int) string {
	if (lexer.currentIndex + len) > lexer.dataLength {
		return ""
	}
	boundIndex := lexer.currentIndex + len
	v := lexer.dataContent[lexer.currentIndex:boundIndex]
	return string(v)
}

func (lexer *schemaLexer) putBack() {
	if lexer.currentIndex > 0 {
		lexer.currentIndex--
	}
}

func (lexer *schemaLexer) fetchText(lval *yySymType, startIndex int) string {
	v := string(lexer.dataContent[startIndex:lexer.currentIndex])
	lval.text = v
	return v
}

func (lexer *schemaLexer) stateTransitQuotedString(lval *yySymType, result []rune, quoteChar, inputChar rune) []rune {
	if inputChar == '\u005C' {
		result = lexer.escapedQuotedCharacter(result)
	} else if inputChar == quoteChar {
		lval.text = string(result)
	} else {
		result = append(result, inputChar)
	}
	return result
}

func (lexer *schemaLexer) escapedQuotedCharacter(result []rune) []rune {
	var escapedCh rune
	v := lexer.peekString(2)
	if (v == "5c") || (v == "5C") {
		escapedCh = '\\'
	} else if v == "27" {
		escapedCh = '\''
	} else if v == "22" {
		escapedCh = '"'
	}
	if escapedCh != 0 {
		result = append(result, escapedCh)
		lexer.currentIndex += 2
	}
	return result
}

func (lexer *schemaLexer) Error(e string) {
	log.Printf("parse error: %s", e)
}
