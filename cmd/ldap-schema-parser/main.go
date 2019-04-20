package main

import (
	"log"
	"os"

	ldapschemaparser "github.com/yinyin/go-ldap-schema-parser"
)

func main() {
	if len(os.Args) < 2 {
		log.Printf("Argument: [LDAP_SCHEMA_TEXT] ...")
	}
	for _, arg := range os.Args[1:] {
		log.Printf("Schema: %v", arg)
		schema, err := ldapschemaparser.Parse(arg)
		if nil != err {
			log.Printf("- ERR: %v", err)
		} else {
			log.Printf("- %v", schema)
		}
	}
}
