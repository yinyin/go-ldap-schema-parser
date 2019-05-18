package ldapschemaparser

import (
	"bufio"
	"errors"
	"io"
	"log"
	"os"
	"sort"
	"strings"
)

const lineFieldSeparator string = ":\t"

const (
	recordTypeLDAPSyntaxSchema       string = "ldap-syntax"
	recordTypeMatchingRuleSchema            = "matching-rule"
	recordTypeMatchingRuleUseSchema         = "matching-rule-use"
	recordTypeAttributeTypeSchema           = "attribute-type"
	recordTypeObjectClassSchema             = "object-class"
	recordTypeDITContentRuleSchema          = "dit-content-rule"
	recordTypeDITStructureRuleSchema        = "dit-structure-rule"
	recordTypeNameFormSchema                = "name-form"
)

func sortedMapKey(m map[string]*GenericSchema) (result []string) {
	result = make([]string, 0, len(m))
	for k := range m {
		result = append(result, k)
	}
	sort.Strings(result)
	return
}

// LDAPSchemaStore is a container of LDAP schemas
type LDAPSchemaStore struct {
	ldapSyntaxSchemaIndex       map[string]*GenericSchema
	matchingRuleSchemaIndex     map[string]*GenericSchema
	matchingRuleNameIndex       map[string]*GenericSchema
	matchingRuleUseSchemaIndex  map[string]*GenericSchema
	attributeTypeSchemaIndex    map[string]*GenericSchema
	attributeTypeNameIndex      map[string]*GenericSchema
	objectClassSchemaIndex      map[string]*GenericSchema
	objectClassNameIndex        map[string]*GenericSchema
	ditContentRuleSchemaIndex   map[string]*GenericSchema
	ditStructureRuleSchemaIndex map[string]*GenericSchema
	nameFormSchemaIndex         map[string]*GenericSchema
}

// NewLDAPSchemaStore create an instance of LDAPSchemaStore
func NewLDAPSchemaStore() (store *LDAPSchemaStore) {
	return &LDAPSchemaStore{
		ldapSyntaxSchemaIndex:       make(map[string]*GenericSchema),
		matchingRuleSchemaIndex:     make(map[string]*GenericSchema),
		matchingRuleNameIndex:       make(map[string]*GenericSchema),
		matchingRuleUseSchemaIndex:  make(map[string]*GenericSchema),
		attributeTypeSchemaIndex:    make(map[string]*GenericSchema),
		attributeTypeNameIndex:      make(map[string]*GenericSchema),
		objectClassSchemaIndex:      make(map[string]*GenericSchema),
		objectClassNameIndex:        make(map[string]*GenericSchema),
		ditContentRuleSchemaIndex:   make(map[string]*GenericSchema),
		ditStructureRuleSchemaIndex: make(map[string]*GenericSchema),
		nameFormSchemaIndex:         make(map[string]*GenericSchema),
	}
}

// AddLDAPSyntaxSchemaText add LDAP syntax schema in text form
func (store *LDAPSchemaStore) AddLDAPSyntaxSchemaText(schemaText string) (err error) {
	genericSchema, err := Parse(schemaText)
	if nil != err {
		return
	}
	ldapSyntaxSchema, err := NewLDAPSyntaxSchemaViaGenericSchema(genericSchema)
	if nil != err {
		return
	}
	existedSchema := store.ldapSyntaxSchemaIndex[ldapSyntaxSchema.NumericOID]
	if nil != existedSchema {
		existedSchema.add(genericSchema)
	} else {
		store.ldapSyntaxSchemaIndex[ldapSyntaxSchema.NumericOID] = genericSchema
	}
	return nil
}

// AddMatchingRuleSchemaText add matching rule schema in text form
func (store *LDAPSchemaStore) AddMatchingRuleSchemaText(schemaText string) (err error) {
	genericSchema, err := Parse(schemaText)
	if nil != err {
		return
	}
	matchingRuleSchema, err := NewMatchingRuleSchemaViaGenericSchema(genericSchema)
	if nil != err {
		return
	}
	existedSchema := store.matchingRuleSchemaIndex[matchingRuleSchema.NumericOID]
	if nil != existedSchema {
		existedSchema.add(genericSchema)
		genericSchema = existedSchema
	} else {
		store.matchingRuleSchemaIndex[matchingRuleSchema.NumericOID] = genericSchema
	}
	for _, name := range matchingRuleSchema.Name {
		if existedSchema = store.matchingRuleNameIndex[name]; nil != existedSchema {
			if existedSchema == genericSchema {
				return nil
			}
			log.Printf("WARN: over-writing matchingRuleNameIndex (name=%v): %v <= %v", name, existedSchema, genericSchema)
		}
		store.matchingRuleNameIndex[name] = genericSchema
	}
	return nil
}

// AddMatchingRuleUseSchemaText add matching rule use schema in text form
func (store *LDAPSchemaStore) AddMatchingRuleUseSchemaText(schemaText string) (err error) {
	genericSchema, err := Parse(schemaText)
	if nil != err {
		return
	}
	matchingRuleUseSchema, err := NewMatchingRuleUseSchemaViaGenericSchema(genericSchema)
	if nil != err {
		return
	}
	existedSchema := store.matchingRuleUseSchemaIndex[matchingRuleUseSchema.NumericOID]
	if nil != existedSchema {
		existedSchema.add(genericSchema)
	} else {
		store.matchingRuleUseSchemaIndex[matchingRuleUseSchema.NumericOID] = genericSchema
	}
	return nil
}

// AddAttributeTypeSchemaText add attribute type schema in text form
func (store *LDAPSchemaStore) AddAttributeTypeSchemaText(schemaText string) (err error) {
	genericSchema, err := Parse(schemaText)
	if nil != err {
		return
	}
	attributeTypeSchema, err := NewAttributeTypeSchemaViaGenericSchema(genericSchema)
	if nil != err {
		return
	}
	existedSchema := store.attributeTypeSchemaIndex[attributeTypeSchema.NumericOID]
	if nil != existedSchema {
		existedSchema.add(genericSchema)
		genericSchema = existedSchema
	} else {
		store.attributeTypeSchemaIndex[attributeTypeSchema.NumericOID] = genericSchema
	}
	for _, name := range attributeTypeSchema.Name {
		if existedSchema = store.attributeTypeNameIndex[name]; nil != existedSchema {
			if existedSchema == genericSchema {
				return nil
			}
			log.Printf("WARN: over-writing attributeTypeNameIndex (name=%v): %v <= %v", name, existedSchema, genericSchema)
		}
		store.attributeTypeNameIndex[name] = genericSchema
	}
	return nil
}

// AddObjectClassSchemaText add object class schema in text form
func (store *LDAPSchemaStore) AddObjectClassSchemaText(schemaText string) (err error) {
	genericSchema, err := Parse(schemaText)
	if nil != err {
		return
	}
	objectClassSchema, err := NewObjectClassSchemaViaGenericSchema(genericSchema)
	if nil != err {
		return
	}
	existedSchema := store.objectClassSchemaIndex[objectClassSchema.NumericOID]
	if nil != existedSchema {
		existedSchema.add(genericSchema)
		genericSchema = existedSchema
	} else {
		store.objectClassSchemaIndex[objectClassSchema.NumericOID] = genericSchema
	}
	for _, name := range objectClassSchema.Name {
		if existedSchema = store.objectClassNameIndex[name]; nil != existedSchema {
			if existedSchema == genericSchema {
				return nil
			}
			log.Printf("WARN: over-writing objectClassNameIndex (name=%v): %v <= %v", name, existedSchema, genericSchema)
		}
		store.objectClassNameIndex[name] = genericSchema
	}
	return nil
}

// AddDITContentRuleSchemaText add DIT content rule schema in text form
func (store *LDAPSchemaStore) AddDITContentRuleSchemaText(schemaText string) (err error) {
	genericSchema, err := Parse(schemaText)
	if nil != err {
		return
	}
	ditContentRuleSchema, err := NewDITContentRuleSchemaViaGenericSchema(genericSchema)
	if nil != err {
		return
	}
	existedSchema := store.ditContentRuleSchemaIndex[ditContentRuleSchema.NumericOID]
	if nil != existedSchema {
		existedSchema.add(genericSchema)
	} else {
		store.ditContentRuleSchemaIndex[ditContentRuleSchema.NumericOID] = genericSchema
	}
	return nil
}

// AddDITStructureRuleSchemaText add DIT structure rule schema in text form
func (store *LDAPSchemaStore) AddDITStructureRuleSchemaText(schemaText string) (err error) {
	genericSchema, err := Parse(schemaText)
	if nil != err {
		return
	}
	ditStructureRuleSchema, err := NewDITStructureRuleSchemaViaGenericSchema(genericSchema)
	if nil != err {
		return
	}
	existedSchema := store.ditStructureRuleSchemaIndex[ditStructureRuleSchema.RuleID]
	if nil != existedSchema {
		existedSchema.add(genericSchema)
	} else {
		store.ditStructureRuleSchemaIndex[ditStructureRuleSchema.RuleID] = genericSchema
	}
	return nil
}

// AddNameFormSchemaText add name form schema in text form
func (store *LDAPSchemaStore) AddNameFormSchemaText(schemaText string) (err error) {
	genericSchema, err := Parse(schemaText)
	if nil != err {
		return
	}
	nameFormSchema, err := NewNameFormSchemaViaGenericSchema(genericSchema)
	if nil != err {
		return
	}
	existedSchema := store.nameFormSchemaIndex[nameFormSchema.NumericOID]
	if nil != existedSchema {
		existedSchema.add(genericSchema)
	} else {
		store.nameFormSchemaIndex[nameFormSchema.NumericOID] = genericSchema
	}
	return nil
}

func (store *LDAPSchemaStore) writeLDAPSyntaxSchema(fp *os.File) (err error) {
	for _, oid := range sortedMapKey(store.ldapSyntaxSchemaIndex) {
		genericSchema := store.ldapSyntaxSchemaIndex[oid]
		ldapSyntaxSchema, err := NewLDAPSyntaxSchemaViaGenericSchema(genericSchema)
		if nil != err {
			log.Printf("ERROR: cannot create LDAP syntax schema object from generic schema [%v]: %v", oid, err)
			continue
		}
		line := recordTypeLDAPSyntaxSchema + lineFieldSeparator + ldapSyntaxSchema.String() + "\n"
		if _, err = fp.WriteString(line); nil != err {
			return err
		}
	}
	return nil
}

func (store *LDAPSchemaStore) writeMatchingRuleSchema(fp *os.File) (err error) {
	for _, oid := range sortedMapKey(store.matchingRuleSchemaIndex) {
		genericSchema := store.matchingRuleSchemaIndex[oid]
		matchingRuleSchema, err := NewMatchingRuleSchemaViaGenericSchema(genericSchema)
		if nil != err {
			log.Printf("ERROR: cannot create matching rule schema object from generic schema [%v]: %v", oid, err)
			continue
		}
		line := recordTypeMatchingRuleSchema + lineFieldSeparator + matchingRuleSchema.String() + "\n"
		if _, err = fp.WriteString(line); nil != err {
			return err
		}
	}
	return nil
}

func (store *LDAPSchemaStore) writeMatchingRuleUseSchema(fp *os.File) (err error) {
	for _, oid := range sortedMapKey(store.matchingRuleUseSchemaIndex) {
		genericSchema := store.matchingRuleSchemaIndex[oid]
		matchingRuleUseSchema, err := NewMatchingRuleUseSchemaViaGenericSchema(genericSchema)
		if nil != err {
			log.Printf("ERROR: cannot create matching rule use schema object from generic schema [%v]: %v", oid, err)
			continue
		}
		line := recordTypeMatchingRuleUseSchema + lineFieldSeparator + matchingRuleUseSchema.String() + "\n"
		if _, err = fp.WriteString(line); nil != err {
			return err
		}
	}
	return nil
}

func (store *LDAPSchemaStore) writeAttributeTypeSchema(fp *os.File) (err error) {
	for _, oid := range sortedMapKey(store.attributeTypeSchemaIndex) {
		genericSchema := store.attributeTypeSchemaIndex[oid]
		attributeTypeSchema, err := NewAttributeTypeSchemaViaGenericSchema(genericSchema)
		if nil != err {
			log.Printf("ERROR: cannot create attribute type schema object from generic schema [%v]: %v", oid, err)
			continue
		}
		line := recordTypeAttributeTypeSchema + lineFieldSeparator + attributeTypeSchema.String() + "\n"
		if _, err = fp.WriteString(line); nil != err {
			return err
		}
	}
	return nil
}

func (store *LDAPSchemaStore) writeObjectClassSchema(fp *os.File) (err error) {
	for _, oid := range sortedMapKey(store.objectClassSchemaIndex) {
		genericSchema := store.objectClassSchemaIndex[oid]
		objectClassSchema, err := NewObjectClassSchemaViaGenericSchema(genericSchema)
		if nil != err {
			log.Printf("ERROR: cannot create object class schema object from generic schema [%v]: %v", oid, err)
			continue
		}
		line := recordTypeObjectClassSchema + lineFieldSeparator + objectClassSchema.String() + "\n"
		if _, err = fp.WriteString(line); nil != err {
			return err
		}
	}
	return nil
}

func (store *LDAPSchemaStore) writeDITContentRuleSchema(fp *os.File) (err error) {
	for _, oid := range sortedMapKey(store.ditContentRuleSchemaIndex) {
		genericSchema := store.objectClassSchemaIndex[oid]
		ditContentRuleSchema, err := NewDITContentRuleSchemaViaGenericSchema(genericSchema)
		if nil != err {
			log.Printf("ERROR: cannot create DIT content rule schema object from generic schema [%v]: %v", oid, err)
			continue
		}
		line := recordTypeDITContentRuleSchema + lineFieldSeparator + ditContentRuleSchema.String() + "\n"
		if _, err = fp.WriteString(line); nil != err {
			return err
		}
	}
	return nil
}

func (store *LDAPSchemaStore) writeDITStructureRuleSchema(fp *os.File) (err error) {
	for _, ruleID := range sortedMapKey(store.ditStructureRuleSchemaIndex) {
		genericSchema := store.objectClassSchemaIndex[ruleID]
		ditStructureRuleSchema, err := NewDITStructureRuleSchemaViaGenericSchema(genericSchema)
		if nil != err {
			log.Printf("ERROR: cannot create DIT structure rule schema object from generic schema [%v]: %v", ruleID, err)
			continue
		}
		line := recordTypeDITStructureRuleSchema + lineFieldSeparator + ditStructureRuleSchema.String() + "\n"
		if _, err = fp.WriteString(line); nil != err {
			return err
		}
	}
	return nil
}

func (store *LDAPSchemaStore) writeNameFormSchema(fp *os.File) (err error) {
	for _, oid := range sortedMapKey(store.nameFormSchemaIndex) {
		genericSchema := store.nameFormSchemaIndex[oid]
		nameFormSchema, err := NewNameFormSchemaViaGenericSchema(genericSchema)
		if nil != err {
			log.Printf("ERROR: cannot create name form schema object from generic schema [%v]: %v", oid, err)
			continue
		}
		line := recordTypeNameFormSchema + lineFieldSeparator + nameFormSchema.String() + "\n"
		if _, err = fp.WriteString(line); nil != err {
			return err
		}
	}
	return nil
}

// WriteToFile write content of store into file at given path
func (store *LDAPSchemaStore) WriteToFile(name string) (err error) {
	fp, err := os.Create(name)
	if nil != err {
		return
	}
	defer fp.Close()
	if err = store.writeLDAPSyntaxSchema(fp); nil != err {
		return
	}
	if err = store.writeMatchingRuleSchema(fp); nil != err {
		return
	}
	if err = store.writeMatchingRuleUseSchema(fp); nil != err {
		return
	}
	if err = store.writeAttributeTypeSchema(fp); nil != err {
		return
	}
	if err = store.writeObjectClassSchema(fp); nil != err {
		return
	}
	if err = store.writeDITContentRuleSchema(fp); nil != err {
		return
	}
	if err = store.writeDITStructureRuleSchema(fp); nil != err {
		return
	}
	if err = store.writeNameFormSchema(fp); nil != err {
		return
	}
	return nil
}

func (store *LDAPSchemaStore) readLine(ln string) (err error) {
	ln = strings.TrimSpace(ln)
	idx := strings.Index(ln, lineFieldSeparator)
	if idx < 0 {
		if len(ln) > 0 {
			log.Printf("WARN: dropping line - [%v]", ln)
		}
		return nil
	}
	k := ln[0:idx]
	v := strings.TrimSpace(ln[idx+len(lineFieldSeparator):])
	switch k {
	case recordTypeLDAPSyntaxSchema:
		err = store.AddLDAPSyntaxSchemaText(v)
	case recordTypeMatchingRuleSchema:
		err = store.AddMatchingRuleSchemaText(v)
	case recordTypeMatchingRuleUseSchema:
		err = store.AddMatchingRuleUseSchemaText(v)
	case recordTypeAttributeTypeSchema:
		err = store.AddAttributeTypeSchemaText(v)
	case recordTypeObjectClassSchema:
		err = store.AddObjectClassSchemaText(v)
	case recordTypeDITContentRuleSchema:
		err = store.AddDITContentRuleSchemaText(v)
	case recordTypeDITStructureRuleSchema:
		err = store.AddDITStructureRuleSchemaText(v)
	case recordTypeNameFormSchema:
		err = store.AddNameFormSchemaText(v)
	default:
		err = errors.New("unknown record type key: " + k)
	}
	return
}

// ReadFromFile read content into store from file at given path
func (store *LDAPSchemaStore) ReadFromFile(name string) (err error) {
	fp, err := os.Open(name)
	if nil != err {
		return
	}
	defer fp.Close()
	reader := bufio.NewReader(fp)
	num := 0
	for {
		ln, err := reader.ReadString('\n')
		num++
		errParse := store.readLine(ln)
		if nil != err {
			if io.EOF == err {
				break
			}
			log.Printf("ERROR: failed on reading from file (file=%v, line=%d, err=%v)", name, num, err)
			return err
		}
		if nil != errParse {
			log.Printf("ERROR: failed on parsing schema text from file (file=%v, line=%d, err=%v)", name, num, errParse)
			return errParse
		}
	}
	return nil
}
