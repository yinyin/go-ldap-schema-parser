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
%token <text> KEYWORD X_KEYWORD NOIDS_ATTR_KEYWORD QSTRINGS_ATTR_KEYWORD
%token <text> SQSTRING DQSTRING

%type <genericSchema> schema attributeDefinitions attributeDefinition
%type <parameterizedKeyword> noids dollarOIDs numberIDS
%type <parameterizedKeyword> qstrings spacedQuotedStrings
%type <text> oid
%type <text> quotedString

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

quotedString: SQSTRING {
  $$ = $1
}
| DQSTRING {
  $$ = $1
}

spacedQuotedStrings: quotedString {
  $$ = newParameterizedKeywordWithParameter($1, QuotedStringsRule)
}
| spacedQuotedStrings SPACES quotedString {
  $$ = $1
  $$.addParameter($3)
}

qstrings: quotedString {
  $$ = newParameterizedKeywordWithParameter($1, QuotedStringsRule)
}
| '(' optionalSpace spacedQuotedStrings optionalSpace ')' {
  $$ = $3
}

attributeDefinition: KEYWORD {
  $$ = newGenericSchema()
  $$.addFlagKeywords($1)
}
| NOIDS_ATTR_KEYWORD optionalSpace noids {
  $$ = newGenericSchema()
  $$.addParameterizedKeyword($1, $3)
}
| QSTRINGS_ATTR_KEYWORD optionalSpace qstrings {
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
