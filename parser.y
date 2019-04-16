%{
package ldapschemaparser
%}

%union{
  text string
}

%token '(' ')' '{' '}'
%token SPACES
%token NUMBER

%token <text> NUMERIC_OID KEYWORD SQSTRING DQSTRING

%start schema

%%

schema: '(' optionalSpace NUMERIC_OID optionalSpace ')'

optionalSpace:
|	SPACES

%%
