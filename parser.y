%{
package ldapschemaparser
%}

%union{
  genericSchema *GenericSchema
  parameterizedKeyword *ParameterizedKeyword
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
%type <parameterizedKeyword> oids dollarOIDs
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
  $$ = newParameterizedKeywordWithParameter($1, OIDsRule)
}
| dollarOIDs optionalSpace '$' optionalSpace oid {
  $$ = $1
  $$.addParameter($5)
}

oids: oid {
  $$ = newParameterizedKeywordWithParameter($1, OIDsRule)
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
  $$.addParameterizedKeyword($1, $3)
}

attributeDefinitions: attributeDefinition {
  $$ = $1
}
| attributeDefinitions SPACES attributeDefinition {
  $1.add($3)
  $$ = $1
}

%%
