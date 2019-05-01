package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

var keywordTypeIdentifierMapping = map[string]string{
	"OIDs":              "NOIDS_ATTR_KEYWORD",
	"RuleIDS":           "NOIDS_ATTR_KEYWORD",
	"OID":               "NOIDS_ATTR_KEYWORD",
	"QuotedDescriptorS": "QSTRINGS_ATTR_KEYWORD",
}

const mappingTrap = "\\s*\"([A-Z]+)\"[`\\s\\\\s+]+\\*\\*([A-Za-z]+)\\*\\*\\s*"

func checkCompatiableTargetRule(keyword, a, b string) bool {
	if keyword == "SUP" {
		if (a == "OIDs" || a == "RuleIDS") && (b == "OIDs" || b == "RuleIDS") {
			return true
		}
	} else if keyword == "FORM" {
		if (a == "OID" || a == "OIDs") && (b == "OID" || b == "OIDs") {
			return true
		}
	}
	if a == b {
		return true
	}
	return false
}

func fetchKeywordTypeMapping(filePath string) (result map[string]string, err error) {
	trap, err := regexp.Compile(mappingTrap)
	if nil != err {
		return
	}
	fp, err := os.Open(filePath)
	if nil != err {
		return
	}
	defer fp.Close()
	reader := bufio.NewReader(fp)
	lineNum := 0
	result = make(map[string]string)
	for {
		line, err := reader.ReadString('\n')
		lineNum++
		m := trap.FindStringSubmatchIndex(line)
		if nil != m {
			keywordText := line[m[2]:m[3]]
			typeText := line[m[4]:m[5]]
			prevType, ok := result[keywordText]
			if ok && !checkCompatiableTargetRule(keywordText, prevType, typeText) {
				err = fmt.Errorf("mis-match with existed type: keyword=%v, type=%v, prev-type=%v, line=%d", keywordText, typeText, prevType, lineNum)
				return nil, err
			}
			result[keywordText] = typeText
		} else {
			line = strings.TrimSpace(line)
			if (line != "") && (line != "`\\s* \")\"`") && (line[0:1] != "#") {
				log.Printf("no-match: [%d] %v", lineNum, line)
			}
		}
		if nil != err {
			log.Printf("stop reading file: %v", err)
			if io.EOF == err {
				err = nil
			}
			break
		}
	}
	return result, err
}

func generateKeywordTypeMap(filePath string, keywordTypeMap map[string]string) (err error) {
	fp, err := os.Create(filePath)
	if nil != err {
		return
	}
	defer fp.Close()
	fp.WriteString("package ldapschemaparser\n\n")
	fp.WriteString("var keywordTypeLookupMap = map[string]int{\n")
	resultCode := make([]string, 0, 0)
	for keywordText, typeText := range keywordTypeMap {
		if typeIdentifier, ok := keywordTypeIdentifierMapping[typeText]; ok {
			aux := strconv.Quote(keywordText) + ": " + typeIdentifier
			resultCode = append(resultCode, aux)
		} else {
			log.Printf("WARN: cannot map typeText to typeIdentifier: %v for %v", typeText, keywordText)
		}
	}
	sort.Strings(resultCode)
	for _, codeLine := range resultCode {
		fp.WriteString("\t" + codeLine + ",\n")
	}
	fp.WriteString("}\n")
	return nil
}

func main() {
	inputFilePath, outputFilePath, err := parseCommandParam()
	if nil != err {
		log.Fatalf("failed on parsing command parameters: %v", err)
		return
	}
	keywordTypeMap, err := fetchKeywordTypeMapping(inputFilePath)
	if nil != err {
		log.Fatalf("failed on fetching keyword type mapping: %v", err)
		return
	}
	log.Printf("fetched keyword type map: %v", keywordTypeMap)
	if err = generateKeywordTypeMap(outputFilePath, keywordTypeMap); nil != err {
		log.Fatalf("failed on generating keyword type map: %v", err)
	}
}
