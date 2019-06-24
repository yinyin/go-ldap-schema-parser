package ldapschemaparser

import (
	"bufio"
	"encoding/json"
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

func (store *LDAPSchemaStore) addLDAPSyntaxGenericSchema(genericSchema *GenericSchema) (err error) {
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

// AddLDAPSyntaxSchemaText add LDAP syntax schema in text form
func (store *LDAPSchemaStore) AddLDAPSyntaxSchemaText(schemaText string) (err error) {
	genericSchema, err := Parse(schemaText)
	if nil != err {
		return
	}
	return store.addLDAPSyntaxGenericSchema(genericSchema)
}

func (store *LDAPSchemaStore) addMatchingRuleGenericSchema(genericSchema *GenericSchema) (err error) {
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
		lowercaseName := strings.ToLower(name)
		if existedSchema = store.matchingRuleNameIndex[lowercaseName]; nil != existedSchema {
			if existedSchema == genericSchema {
				return nil
			}
			log.Printf("WARN: over-writing matchingRuleNameIndex (name=%v): %v <= %v", name, existedSchema, genericSchema)
		}
		store.matchingRuleNameIndex[lowercaseName] = genericSchema
	}
	return nil
}

// AddMatchingRuleSchemaText add matching rule schema in text form
func (store *LDAPSchemaStore) AddMatchingRuleSchemaText(schemaText string) (err error) {
	genericSchema, err := Parse(schemaText)
	if nil != err {
		return
	}
	return store.addMatchingRuleGenericSchema(genericSchema)
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

func (store *LDAPSchemaStore) addAttributeTypeGenericSchema(genericSchema *GenericSchema) (err error) {
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
		lowercaseName := strings.ToLower(name)
		if existedSchema = store.attributeTypeNameIndex[lowercaseName]; nil != existedSchema {
			if existedSchema == genericSchema {
				return nil
			}
			log.Printf("WARN: over-writing attributeTypeNameIndex (name=%v): %v <= %v", name, existedSchema, genericSchema)
		}
		store.attributeTypeNameIndex[lowercaseName] = genericSchema
	}
	return nil
}

// AddAttributeTypeSchemaText add attribute type schema in text form
func (store *LDAPSchemaStore) AddAttributeTypeSchemaText(schemaText string) (err error) {
	genericSchema, err := Parse(schemaText)
	if nil != err {
		return
	}
	return store.addAttributeTypeGenericSchema(genericSchema)
}

func (store *LDAPSchemaStore) addObjectClassGenericSchema(genericSchema *GenericSchema) (err error) {
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
		lowercaseName := strings.ToLower(name)
		if existedSchema = store.objectClassNameIndex[lowercaseName]; nil != existedSchema {
			if existedSchema == genericSchema {
				return nil
			}
			log.Printf("WARN: over-writing objectClassNameIndex (name=%v): %v <= %v", name, existedSchema, genericSchema)
		}
		store.objectClassNameIndex[lowercaseName] = genericSchema
	}
	return nil
}

// AddObjectClassSchemaText add object class schema in text form
func (store *LDAPSchemaStore) AddObjectClassSchemaText(schemaText string) (err error) {
	genericSchema, err := Parse(schemaText)
	if nil != err {
		return
	}
	return store.addObjectClassGenericSchema(genericSchema)
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

func (store *LDAPSchemaStore) writeFieldSeparatedSchemaTexts(fp *os.File, recordType string, schemaTexts []string) (err error) {
	for _, text := range schemaTexts {
		line := recordType + lineFieldSeparator + text + "\n"
		if _, err = fp.WriteString(line); nil != err {
			return err
		}
	}
	return nil
}

func (store *LDAPSchemaStore) collectLDAPSyntaxSchemaTexts(stopOnError bool) (result []string, err error) {
	for _, oid := range sortedMapKey(store.ldapSyntaxSchemaIndex) {
		genericSchema := store.ldapSyntaxSchemaIndex[oid]
		ldapSyntaxSchema, err := NewLDAPSyntaxSchemaViaGenericSchema(genericSchema)
		if nil != err {
			log.Printf("ERROR: cannot create LDAP syntax schema object from generic schema [%v]: %v", oid, err)
			if stopOnError {
				return nil, err
			}
			continue
		}
		result = append(result, ldapSyntaxSchema.String())
	}
	return result, nil
}

func (store *LDAPSchemaStore) writeLDAPSyntaxSchema(fp *os.File) (err error) {
	schemaTexts, err := store.collectLDAPSyntaxSchemaTexts(false)
	if nil != err {
		return
	}
	return store.writeFieldSeparatedSchemaTexts(fp, recordTypeLDAPSyntaxSchema, schemaTexts)
}

func (store *LDAPSchemaStore) collectMatchingRuleSchemaTexts(stopOnError bool) (result []string, err error) {
	for _, oid := range sortedMapKey(store.matchingRuleSchemaIndex) {
		genericSchema := store.matchingRuleSchemaIndex[oid]
		matchingRuleSchema, err := NewMatchingRuleSchemaViaGenericSchema(genericSchema)
		if nil != err {
			log.Printf("ERROR: cannot create matching rule schema object from generic schema [%v]: %v", oid, err)
			if stopOnError {
				return nil, err
			}
			continue
		}
		result = append(result, matchingRuleSchema.String())
	}
	return result, nil
}

func (store *LDAPSchemaStore) writeMatchingRuleSchema(fp *os.File) (err error) {
	schemaTexts, err := store.collectMatchingRuleSchemaTexts(false)
	if nil != err {
		return
	}
	return store.writeFieldSeparatedSchemaTexts(fp, recordTypeMatchingRuleSchema, schemaTexts)
}

func (store *LDAPSchemaStore) collectMatchingRuleUseSchemaTexts(stopOnError bool) (result []string, err error) {
	for _, oid := range sortedMapKey(store.matchingRuleUseSchemaIndex) {
		genericSchema := store.matchingRuleSchemaIndex[oid]
		matchingRuleUseSchema, err := NewMatchingRuleUseSchemaViaGenericSchema(genericSchema)
		if nil != err {
			log.Printf("ERROR: cannot create matching rule use schema object from generic schema [%v]: %v", oid, err)
			if stopOnError {
				return nil, err
			}
			continue
		}
		result = append(result, matchingRuleUseSchema.String())
	}
	return result, nil
}

func (store *LDAPSchemaStore) writeMatchingRuleUseSchema(fp *os.File) (err error) {
	schemaTexts, err := store.collectMatchingRuleUseSchemaTexts(false)
	if nil != err {
		return
	}
	return store.writeFieldSeparatedSchemaTexts(fp, recordTypeMatchingRuleUseSchema, schemaTexts)
}

func (store *LDAPSchemaStore) collectAttributeTypeSchemaTexts(stopOnError bool) (result []string, err error) {
	for _, oid := range sortedMapKey(store.attributeTypeSchemaIndex) {
		genericSchema := store.attributeTypeSchemaIndex[oid]
		attributeTypeSchema, err := NewAttributeTypeSchemaViaGenericSchema(genericSchema)
		if nil != err {
			log.Printf("ERROR: cannot create attribute type schema object from generic schema [%v]: %v", oid, err)
			if stopOnError {
				return nil, err
			}
			continue
		}
		result = append(result, attributeTypeSchema.String())
	}
	return result, nil
}

func (store *LDAPSchemaStore) writeAttributeTypeSchema(fp *os.File) (err error) {
	schemaTexts, err := store.collectAttributeTypeSchemaTexts(false)
	if nil != err {
		return
	}
	return store.writeFieldSeparatedSchemaTexts(fp, recordTypeAttributeTypeSchema, schemaTexts)
}

func (store *LDAPSchemaStore) collectObjectClassSchemaTexts(stopOnError bool) (result []string, err error) {
	for _, oid := range sortedMapKey(store.objectClassSchemaIndex) {
		genericSchema := store.objectClassSchemaIndex[oid]
		objectClassSchema, err := NewObjectClassSchemaViaGenericSchema(genericSchema)
		if nil != err {
			log.Printf("ERROR: cannot create object class schema object from generic schema [%v]: %v", oid, err)
			if stopOnError {
				return nil, err
			}
			continue
		}
		result = append(result, objectClassSchema.String())
	}
	return result, nil
}

func (store *LDAPSchemaStore) writeObjectClassSchema(fp *os.File) (err error) {
	schemaTexts, err := store.collectObjectClassSchemaTexts(false)
	if nil != err {
		return
	}
	return store.writeFieldSeparatedSchemaTexts(fp, recordTypeObjectClassSchema, schemaTexts)
}

func (store *LDAPSchemaStore) collectDITContentRuleSchemaTexts(stopOnError bool) (result []string, err error) {
	for _, oid := range sortedMapKey(store.ditContentRuleSchemaIndex) {
		genericSchema := store.objectClassSchemaIndex[oid]
		ditContentRuleSchema, err := NewDITContentRuleSchemaViaGenericSchema(genericSchema)
		if nil != err {
			log.Printf("ERROR: cannot create DIT content rule schema object from generic schema [%v]: %v", oid, err)
			if stopOnError {
				return nil, err
			}
			continue
		}
		result = append(result, ditContentRuleSchema.String())
	}
	return result, nil
}

func (store *LDAPSchemaStore) writeDITContentRuleSchema(fp *os.File) (err error) {
	schemaTexts, err := store.collectDITContentRuleSchemaTexts(false)
	if nil != err {
		return
	}
	return store.writeFieldSeparatedSchemaTexts(fp, recordTypeDITContentRuleSchema, schemaTexts)
}

func (store *LDAPSchemaStore) collectDITStructureRuleSchemaTexts(stopOnError bool) (result []string, err error) {
	for _, ruleID := range sortedMapKey(store.ditStructureRuleSchemaIndex) {
		genericSchema := store.objectClassSchemaIndex[ruleID]
		ditStructureRuleSchema, err := NewDITStructureRuleSchemaViaGenericSchema(genericSchema)
		if nil != err {
			log.Printf("ERROR: cannot create DIT structure rule schema object from generic schema [%v]: %v", ruleID, err)
			if stopOnError {
				return nil, err
			}
			continue
		}
		result = append(result, ditStructureRuleSchema.String())
	}
	return result, nil
}

func (store *LDAPSchemaStore) writeDITStructureRuleSchema(fp *os.File) (err error) {
	schemaTexts, err := store.collectDITStructureRuleSchemaTexts(false)
	if nil != err {
		return
	}
	return store.writeFieldSeparatedSchemaTexts(fp, recordTypeDITStructureRuleSchema, schemaTexts)
}

func (store *LDAPSchemaStore) collectNameFormSchemaTexts(stopOnError bool) (result []string, err error) {
	for _, oid := range sortedMapKey(store.nameFormSchemaIndex) {
		genericSchema := store.nameFormSchemaIndex[oid]
		nameFormSchema, err := NewNameFormSchemaViaGenericSchema(genericSchema)
		if nil != err {
			log.Printf("ERROR: cannot create name form schema object from generic schema [%v]: %v", oid, err)
			if stopOnError {
				return nil, err
			}
			continue
		}
		result = append(result, nameFormSchema.String())
	}
	return result, nil
}

func (store *LDAPSchemaStore) writeNameFormSchema(fp *os.File) (err error) {
	schemaTexts, err := store.collectNameFormSchemaTexts(false)
	if nil != err {
		return
	}
	return store.writeFieldSeparatedSchemaTexts(fp, recordTypeNameFormSchema, schemaTexts)
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

// WriteToJSONFile write content of store into given path in JSON form.
func (store *LDAPSchemaStore) WriteToJSONFile(name string) (err error) {
	var aux struct {
		LDAPSyntax       []string `json:"ldap_syntax,omitempty"`
		MatchingRule     []string `json:"matching_rule,omitempty"`
		MatchingRuleUse  []string `json:"matching_rule_use,omitempty"`
		AttributeType    []string `json:"attribute_type,omitempty"`
		ObjectClass      []string `json:"object_class,omitempty"`
		DITContentRule   []string `json:"dit_content_rule,omitempty"`
		DITStructureRule []string `json:"dit_structure_rule,omitempty"`
		NameForm         []string `json:"name_form,omitempty"`
	}
	if aux.LDAPSyntax, err = store.collectLDAPSyntaxSchemaTexts(true); nil != err {
		return
	}
	if aux.MatchingRule, err = store.collectMatchingRuleSchemaTexts(true); nil != err {
		return
	}
	if aux.MatchingRuleUse, err = store.collectMatchingRuleUseSchemaTexts(true); nil != err {
		return
	}
	if aux.AttributeType, err = store.collectAttributeTypeSchemaTexts(true); nil != err {
		return
	}
	if aux.ObjectClass, err = store.collectObjectClassSchemaTexts(true); nil != err {
		return
	}
	if aux.DITContentRule, err = store.collectDITContentRuleSchemaTexts(true); nil != err {
		return
	}
	if aux.DITStructureRule, err = store.collectDITStructureRuleSchemaTexts(true); nil != err {
		return
	}
	if aux.NameForm, err = store.collectNameFormSchemaTexts(true); nil != err {
		return
	}
	buf, err := json.Marshal(&aux)
	if nil != err {
		return
	}
	fp, err := os.Create(name)
	if nil != err {
		return
	}
	defer fp.Close()
	_, err = fp.Write(buf)
	return err
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

func (store *LDAPSchemaStore) pullObjectClassWhenNotExist(source *LDAPSchemaStore, verbose bool, dependentRefName string, objectClassName string) (err error) {
	if _, ok := store.objectClassNameIndex[objectClassName]; ok {
		if verbose {
			log.Printf("INFO: reach object class for %s via name: %s", dependentRefName, objectClassName)
		}
		return nil
	}
	if remoteGenericSchema, ok := source.objectClassNameIndex[objectClassName]; ok {
		if err = store.addObjectClassGenericSchema(remoteGenericSchema); nil != err {
			log.Printf("ERROR: failed on adding dependent object class schema %s for %s from source: %v", objectClassName, dependentRefName, err)
			return err
		} else if verbose {
			log.Printf("INFO: reached object class for %s via name at remote store: %s", dependentRefName, objectClassName)
		}
		return nil
	}
	if verbose {
		log.Printf("ERROR: failed on reach object class for %s: %v", dependentRefName, objectClassName)
	}
	return errors.New("needed object class for " + dependentRefName + " not found: " + objectClassName)
}

func (store *LDAPSchemaStore) pullAttributeTypeWhenNotExist(source *LDAPSchemaStore, verbose bool, dependentRefName string, attributeTypeName string) (err error) {
	lowercaseAttributeTypeName := strings.ToLower(attributeTypeName)
	if _, ok := store.attributeTypeNameIndex[lowercaseAttributeTypeName]; ok {
		if verbose {
			log.Printf("INFO: reach attribute type for %s via name: %s", dependentRefName, attributeTypeName)
		}
		return nil
	}
	if remoteGenericSchema, ok := source.attributeTypeNameIndex[lowercaseAttributeTypeName]; ok {
		if err = store.addAttributeTypeGenericSchema(remoteGenericSchema); nil != err {
			log.Printf("ERROR: failed on adding dependent attribute type schema %s for %s from source: %v", attributeTypeName, dependentRefName, err)
			return err
		} else if verbose {
			log.Printf("INFO: reach attribute type for %s via name at remote store: %s", dependentRefName, attributeTypeName)
		}
		return nil
	}
	if verbose {
		log.Printf("ERROR: failed on reach attribute type for %s: %v", dependentRefName, attributeTypeName)
	}
	return errors.New("needed attribute type for " + dependentRefName + " not found: " + attributeTypeName)
}

func (store *LDAPSchemaStore) pullMatchingRuleWhenNotExist(source *LDAPSchemaStore, verbose bool, dependentRefName string, matchingRuleName string) (err error) {
	lowercaseMatchingRuleName := strings.ToLower(matchingRuleName)
	if _, ok := store.matchingRuleNameIndex[lowercaseMatchingRuleName]; ok {
		if verbose {
			log.Printf("INFO: reach matching rule for %s via name: %s", dependentRefName, matchingRuleName)
		}
		return nil
	}
	if remoteGenericSchema, ok := source.matchingRuleNameIndex[lowercaseMatchingRuleName]; ok {
		if err = store.addMatchingRuleGenericSchema(remoteGenericSchema); nil != err {
			log.Printf("ERROR: failed on adding dependent matching rule schema %s for %s from source: %v", matchingRuleName, dependentRefName, err)
			return err
		} else if verbose {
			log.Printf("INFO: reach matching rule for %s via name at remote store: %s", dependentRefName, matchingRuleName)
		}
		return nil
	}
	if verbose {
		log.Printf("ERROR: failed on reach matching rule for %s: %v", dependentRefName, matchingRuleName)
	}
	return errors.New("needed matching rule for " + dependentRefName + " not found: " + matchingRuleName)
}

func (store *LDAPSchemaStore) pullLDAPSyntaxWhenNotExist(source *LDAPSchemaStore, verbose bool, dependentRefName string, ldapSyntaxOID string) (err error) {
	if _, ok := store.ldapSyntaxSchemaIndex[ldapSyntaxOID]; ok {
		if verbose {
			log.Printf("INFO: reach LDAP syntax for %s via name: %s", dependentRefName, ldapSyntaxOID)
		}
		return nil
	}
	if remoteGenericSchema, ok := source.ldapSyntaxSchemaIndex[ldapSyntaxOID]; ok {
		if err = store.addLDAPSyntaxGenericSchema(remoteGenericSchema); nil != err {
			log.Printf("ERROR: failed on adding dependent LDAP syntax schema %s for %s from source: %v", ldapSyntaxOID, dependentRefName, err)
			return err
		} else if verbose {
			log.Printf("INFO: reach LDAP syntax for %s via name at remote store: %s", dependentRefName, ldapSyntaxOID)
		}
		return nil
	}
	if verbose {
		log.Printf("ERROR: failed on reach LDAP syntax for %s: %v", dependentRefName, ldapSyntaxOID)
	}
	return errors.New("needed LDAP syntax for " + dependentRefName + " not found: " + ldapSyntaxOID)
}

func (store *LDAPSchemaStore) pullObjectClassesDependencies(source *LDAPSchemaStore, verbose bool) (err error) {
	for _, oid := range sortedMapKey(store.objectClassSchemaIndex) {
		genericSchema := store.objectClassSchemaIndex[oid]
		objectClassSchema, err := NewObjectClassSchemaViaGenericSchema(genericSchema)
		if nil != err {
			log.Printf("ERROR: cannot create object class schema object from generic schema for pull dependent schema [%v]: %v", oid, err)
			return err
		}
		for _, superClassName := range objectClassSchema.SuperClasses {
			if err = store.pullObjectClassWhenNotExist(source, verbose, objectClassSchema.NumericOID, superClassName); nil != err {
				return err
			}
		}
		for _, attributeName := range objectClassSchema.Must {
			if err = store.pullAttributeTypeWhenNotExist(source, verbose, objectClassSchema.NumericOID, attributeName); nil != err {
				return err
			}
		}
		for _, attributeName := range objectClassSchema.May {
			if err = store.pullAttributeTypeWhenNotExist(source, verbose, objectClassSchema.NumericOID, attributeName); nil != err {
				return err
			}
		}
	}
	return nil
}

func (store *LDAPSchemaStore) pullAttributeTypesDependencies(source *LDAPSchemaStore, verbose bool) (err error) {
	previousCount := 0
	passCount := 1
	for len(store.attributeTypeSchemaIndex) != previousCount {
		previousCount = len(store.attributeTypeSchemaIndex)
		for _, oid := range sortedMapKey(store.attributeTypeSchemaIndex) {
			genericSchema := store.attributeTypeSchemaIndex[oid]
			attributeTypeSchema, err := NewAttributeTypeSchemaViaGenericSchema(genericSchema)
			if nil != err {
				log.Printf("ERROR: cannot create attribute type schema object from generic schema for pull dependent schema [%v]: %v", oid, err)
				return err
			}
			if "" != attributeTypeSchema.SuperType {
				if err = store.pullAttributeTypeWhenNotExist(source, verbose, attributeTypeSchema.NumericOID, attributeTypeSchema.SuperType); nil != err {
					return err
				}
			}
			if "" != attributeTypeSchema.Equality {
				if err = store.pullMatchingRuleWhenNotExist(source, verbose, attributeTypeSchema.NumericOID, attributeTypeSchema.Equality); nil != err {
					return err
				}
			}
			if "" != attributeTypeSchema.Ordering {
				if err = store.pullMatchingRuleWhenNotExist(source, verbose, attributeTypeSchema.NumericOID, attributeTypeSchema.Ordering); nil != err {
					return err
				}
			}
			if "" != attributeTypeSchema.SubString {
				if err = store.pullMatchingRuleWhenNotExist(source, verbose, attributeTypeSchema.NumericOID, attributeTypeSchema.SubString); nil != err {
					return err
				}
			}
			if "" != attributeTypeSchema.SyntaxOID {
				if err = store.pullLDAPSyntaxWhenNotExist(source, verbose, attributeTypeSchema.NumericOID, attributeTypeSchema.SyntaxOID); nil != err {
					return err
				}
			}
		}
		if verbose {
			log.Printf("INFO: %d pass of pulling attribute type dependencies", passCount)
		}
		passCount = passCount + 1
	}
	return nil
}

func (store *LDAPSchemaStore) pullMatchingRulesDependencies(source *LDAPSchemaStore, verbose bool) (err error) {
	for _, oid := range sortedMapKey(store.matchingRuleSchemaIndex) {
		genericSchema := store.matchingRuleSchemaIndex[oid]
		matchingRuleSchema, err := NewMatchingRuleSchemaViaGenericSchema(genericSchema)
		if nil != err {
			log.Printf("ERROR: cannot create matching rule schema object from generic schema for pull dependent schema [%v]: %v", oid, err)
			return err
		}
		if "" != matchingRuleSchema.Syntax {
			if err = store.pullLDAPSyntaxWhenNotExist(source, verbose, matchingRuleSchema.NumericOID, matchingRuleSchema.Syntax); nil != err {
				return err
			}
		}
	}
	return nil
}

// PullDependentSchema pull schemas used by contained schemas from source store into this store.
func (store *LDAPSchemaStore) PullDependentSchema(source *LDAPSchemaStore, verbose bool) (err error) {
	if err = store.pullObjectClassesDependencies(source, verbose); nil != err {
		log.Printf("ERROR: failed on pull dependecies for object classes: %v", err)
		return
	}
	if err = store.pullAttributeTypesDependencies(source, verbose); nil != err {
		log.Printf("ERROR: failed on pull dependecies for attribute types: %v", err)
		return
	}
	if err = store.pullMatchingRulesDependencies(source, verbose); nil != err {
		log.Printf("ERROR: failed on pull dependecies for matching rules: %v", err)
		return
	}
	return nil
}
