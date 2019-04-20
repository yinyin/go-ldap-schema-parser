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

%token <text> NUMBER
%token <text> NUMERIC_OID
%token <text> KEYWORD X_KEYWORD NOIDS_ATTR_KEYWORD
%token <text> SQSTRING DQSTRING

// %type <text> quotedString
%type <genericSchema> schema attributeDefinitions attributeDefinition
%type <parameterizedKeyword> noids dollarOIDs numberIDS
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

numberIDS: NUMBER {
  $$ = newParameterizedKeywordWithParameter($1, NumberIDsRule)
}
| numberIDS SPACES NUMBER {
  $$ = $1
  $$.addParameter($3)
}

noids: oid {
  $$ = newParameterizedKeywordWithParameter($1, OIDsRule)
}
| '(' optionalSpace dollarOIDs optionalSpace ')' {
  $$ = $3
}
| NUMBER {
  $$ = newParameterizedKeywordWithParameter($1, NumberIDsRule)
}
| '(' optionalSpace numberIDS optionalSpace ')' {
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
| NOIDS_ATTR_KEYWORD optionalSpace noids {
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
