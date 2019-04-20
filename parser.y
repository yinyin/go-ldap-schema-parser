%{
package ldapschemaparser
%}

%union{
  genericSchema *GenericSchema
  text string
}

%token '(' ')' '{' '}'
%token SPACES
%token NUMBER

%token <text> NUMERIC_OID KEYWORD SQSTRING DQSTRING

%type <genericSchema> schema

%start schema

%%

schema: '(' optionalSpace NUMERIC_OID optionalSpace ')' {
  $$ = &GenericSchema {
    NumericOID: $3,
  }
  yylex.(*schemaLexer).result = $$
}

optionalSpace:
|	SPACES

%%
