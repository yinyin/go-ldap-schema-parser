package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"regexp"
	"strings"
	"unicode"
)

const (
	readModeNormal int = iota
	readModeSchema
)

// LineTypePlainText, LineTypeSchema and LineTypeChapter indicate type of line
const (
	LineTypePlainText int = iota
	LineTypeSchema
	LineTypeChapter
)

var trapPageHeader1, trapPageHeader2, trapPageFooter, trapChapter *regexp.Regexp

func init() {
	trapPageHeader1 = regexp.MustCompile("^RFC\\s*[0-9]{4}\\s+")
	trapPageHeader2 = regexp.MustCompile("[A-Z][a-z]{2,8}\\s[0-9]{4}$")
	trapPageFooter = regexp.MustCompile("\\s+\\[Page\\s+[0-9]+\\]$")
	trapChapter = regexp.MustCompile("^(([0-9]+\\.){1,5})\\s+")
}

func checkSchemaStartLine(l string, offset int) bool {
	spaceCount := 0
	numberCount := 0
	dotCount := 0
	othersCount := 0
	for idx, ch := range []rune(l) {
		if idx < (offset + 1) {
			continue
		}
		if ch == ')' {
			break
		} else if ch == ' ' {
			spaceCount++
		} else if (ch >= '0') && (ch <= '9') {
			numberCount++
		} else if ch == '.' {
			dotCount++
		} else {
			othersCount++
		}
	}
	if ((othersCount > 0) && (spaceCount < 2)) || (numberCount < 1) {
		return false
	}
	return true
}

func isSchemaStart(l string) (spaceCount int) {
	for _, ch := range []rune(l) {
		if ' ' == ch {
			spaceCount++
		} else if '(' == ch {
			if !checkSchemaStartLine(l, spaceCount) {
				return -1
			}
			return
		} else {
			return -1
		}
	}
	return -1
}

func countLeadingSpace(l string) (spaceCount int) {
	for _, ch := range []rune(l) {
		if ' ' == ch {
			spaceCount++
		} else {
			break
		}
	}
	return
}

func isSchemaEnd(l string) bool {
	a := []rune(l)
	for idx := len(a) - 1; idx >= 0; idx-- {
		ch := a[idx]
		if ')' == ch {
			return true
		} else if ' ' != ch {
			break
		}
	}
	return false
}

// RFCTextReader is a reader loads strings from RFC text file
// with special handling of LDAP schema text
type RFCTextReader struct {
	fp     *os.File
	reader *bufio.Reader
	mode   int
	lineno int

	lastLine bool

	schemaTextSpaceCount int
	schemaTextBuffer     string

	CurrentChapter string
}

// OpenRFCTextReader open an instance of RFCTextReader
func OpenRFCTextReader(name string) (b *RFCTextReader, err error) {
	fp, err := os.Open(name)
	if nil != err {
		return
	}
	reader := bufio.NewReader(fp)
	b = &RFCTextReader{
		fp:       fp,
		reader:   reader,
		mode:     readModeNormal,
		lastLine: false,
	}
	return b, nil
}

// Close opened file pointer
func (b *RFCTextReader) Close() (err error) {
	return b.fp.Close()
}

// ReadLine get one line from reader
func (b *RFCTextReader) ReadLine() (v string, lineType int, err error) {
	if b.lastLine {
		err = io.EOF
		return
	}
	for !b.lastLine {
		if v, err = b.reader.ReadString('\n'); nil != err {
			if err == io.EOF {
				b.lastLine = true
				err = nil
			} else {
				return
			}
		}
		b.lineno++
		v = strings.TrimRightFunc(v, unicode.IsSpace)
		switch b.mode {
		case readModeNormal:
			if spaceCount := isSchemaStart(v); (spaceCount > 3) && (spaceCount < 12) {
				v = strings.TrimLeftFunc(v, unicode.IsSpace)
				if isSchemaEnd(v) {
					return v, LineTypeSchema, nil
				}
				b.schemaTextSpaceCount = spaceCount
				b.mode = readModeSchema
				b.schemaTextBuffer = v
			} else {
				if chapterText := trapChapter.FindString(v); chapterText != "" {
					b.CurrentChapter = strings.TrimRightFunc(chapterText, unicode.IsSpace)
					lineType = LineTypeChapter
				}
				return
			}
		case readModeSchema:
			if 0 == len(v) {
				break
			}
			if trapPageHeader1.MatchString(v) && trapPageHeader2.MatchString(v) {
				break
			}
			if trapPageFooter.MatchString(v) {
				break
			}
			spaceCount := countLeadingSpace(v)
			if spaceCount < b.schemaTextSpaceCount {
				log.Printf("WARN: indent ot enough for schema: %v, line=%d", v, b.lineno)
			}
			v = strings.TrimLeftFunc(v, unicode.IsSpace)
			b.schemaTextBuffer = b.schemaTextBuffer + " " + v
			if isSchemaEnd(v) {
				b.mode = readModeNormal
				v = b.schemaTextBuffer
				lineType = LineTypeSchema
				return
			}
		default:
			log.Printf("ERR: unknown mode: %v", b.mode)
		}
	}
	return
}
