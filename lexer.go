package ldapschemaparser

//go:generate goyacc -o parser.go parser.y

import (
	"log"
	"unicode"
)

const dataEOF = 0

type ldapSchemaLexer struct {
	dataContent  []rune
	dataLength   int
	currentIndex int
}

func (lexer *ldapSchemaLexer) Lex(lval *yySymType) (lexIdentifier int) {
	var result []rune
	startIndex := lexer.currentIndex
	for {
		ch := lexer.next()
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
				return
			}
		case SQSTRING:
			if ch == '\u005C' {
				v := lexer.peekString(2)
				if (v == "5c") || (v == "5C") {
					result = append(result, '\\')
					lexer.currentIndex += 2
				} else if v == "27" {
					result = append(result, '\'')
					lexer.currentIndex += 2
				}
			} else if ch == '\'' {
				lval.text = string(result)
			} else {
				result = append(result, ch)
			}
		case DQSTRING:
			if ch == '\u005C' {
				v := lexer.peekString(2)
				if (v == "5c") || (v == "5C") {
					result = append(result, '\\')
					lexer.currentIndex += 2
				} else if v == "27" {
					result = append(result, '\'')
					lexer.currentIndex += 2
				} else if v == "22" {
					result = append(result, '"')
					lexer.currentIndex += 2
				}
			} else if ch == '"' {
				lval.text = string(result)
			} else {
				result = append(result, ch)
			}
		}
	}
	return 0
}

func (lexer *ldapSchemaLexer) next() rune {
	if lexer.currentIndex >= lexer.dataLength {
		return dataEOF
	}
	ch := lexer.dataContent[lexer.currentIndex]
	lexer.currentIndex++
	return ch
}

func (lexer *ldapSchemaLexer) peekString(len int) string {
	if (lexer.currentIndex + len) > lexer.dataLength {
		return ""
	}
	boundIndex := lexer.currentIndex + len
	v := lexer.dataContent[lexer.currentIndex:boundIndex]
	return string(v)
}

func (lexer *ldapSchemaLexer) putBack() {
	if lexer.currentIndex > 0 {
		lexer.currentIndex--
	}
}

func (lexer *ldapSchemaLexer) fetchText(lval *yySymType, startIndex int) string {
	v := string(lexer.dataContent[startIndex:lexer.currentIndex])
	lval.text = v
	return v
}

func (lexer *ldapSchemaLexer) Error(e string) {
	log.Printf("parse error: %s", e)
}
