package ldapschemaparser

import (
	"log"
)

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
