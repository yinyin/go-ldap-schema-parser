%{
package ldapschemaparser
%}

%union{
  genericSchema *GenericSchema
  attrWithOIDs *AttributeWithOIDs
  text string
}

%token '(' ')' '{' '}'
%token SPACES
%token NUMBER

%token <text> NUMERIC_OID
%token <text> KEYWORD X_KEYWORD OIDS_ATTR_KEYWORD
%token <text> SQSTRING DQSTRING

// %type <text> quotedString
%type <genericSchema> schema attributeDefinitions attributeDefinition
%type <attrWithOIDs> oids dollarOIDs
%type <text> oid

%start schema

%%

schema: '(' optionalSpace NUMERIC_OID SPACES attributeDefinitions optionalSpace ')' {
  $$ = $5
  $$.NumericOID = $3
  yylex.(*schemaLexer).result = $$
}

optionalSpace:
|	SPACES

oid: NUMERIC_OID {
  $$ = $1
}
| KEYWORD {
  $$ = $1
}

dollarOIDs: oid {
  $$ = newAttributeWithOIDsWithOID($1)
}
| dollarOIDs optionalSpace '$' optionalSpace oid {
  $$ = $1
  $$.addOID($5)
}

oids: oid {
  $$ = newAttributeWithOIDsWithOID($1)
}
| '(' optionalSpace dollarOIDs optionalSpace ')' {
  $$ = $3
}

/*
quotedString: SQSTRING {
  $$ = $1
}
| DQSTRING {
  $$ = $1
}
*/

attributeDefinition: KEYWORD {
  $$ = newGenericSchema()
  $$.addFlagKeywords($1)
}
| OIDS_ATTR_KEYWORD optionalSpace oids {
  $$ = newGenericSchema()
  $$.addAttributeWithOIDs($1, $3)
}

attributeDefinitions: attributeDefinition {
  $$ = $1
}
| attributeDefinitions SPACES attributeDefinition {
  $1.add($3)
  $$ = $1
}

%%
