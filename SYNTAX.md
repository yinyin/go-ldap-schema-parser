Here is syntax of LDAP schema in regular expression for understanding.

# Common Production

* **NumericOID** = `[0-9]+` `( \.[0-9]+ )*`
* **Descriptor** = `[a-zA-Z]` `[a-zA-Z0-9-]*`
* **OID** = **NumericOID** | **Descriptor**

* **NumericOIDwithLength** = `[0-9]+` `( \.[0-9]+ )*` `"{" [0-9]+ "}"`

# Schema Definition

* **OIDs** =
    - **OID** |
    - `"(" \s*` **OID** `( \s* "$" \s* ` **OID** ` )*` `\s* ")"`

* **QuotedDescriptor** = `"'"` **Descriptor** `"'"`
* **QuotedDescriptorS** =
    - **QuotedDescriptor** |
    - `"(" \s*` **QuotedDescriptor** `( \s+ ` **QuotedDescriptor** ` )* ` `\s* ")"`

* **QuotedDString** = `"'" ` `( "\27" | "\5c" | "\5C" | [^'\] )+` ` "'"`
* **QuotedDStringS** =
    - **QuotedDString** |
    - `"(" \s*` **QuotedDString** `( \s+ ` **QuotedDString** ` )* ` `\s* ")"`

* **Extensions** = `"X-" [a-zA-Z-_]+ \s+ ` **QuotedDStringS**

## Object Class

Default *kind* of object class is `"STRUCTURAL"` type.

* **ObjectClassDescription** = `"(" \s*` **NumericOID**
  `( \s+ "NAME"` `\s+` **QuotedDescriptorS** ` )?`
  `( \s+ "DESC"` `\s+` **QuotedDString** ` )?`
  `( \s+ "OBSOLETE" )?`
  `( \s+ "SUP"` `\s+` **OIDs** ` )?`
  `( \s+ ( "ABSTRACT" | "STRUCTURAL" | "AUXILIARY" ) )?`
  `( \s+ "MUST"` `\s+` **OIDs** ` )?`
  `( \s+ "MAY"` `\s+` **OIDs** ` )?`
  `( \s+ ` **Extensions** ` )*`
  `\s* ")"`

## Attribute Type

* **AttributeTypeDescription** = `"(" \s*` **NumericOID**
  `( \s+ "NAME"` `\s+` **QuotedDescriptorS** ` )?`
  `( \s+ "DESC"` `\s+` **QuotedDString** ` )?`
  `( \s+ "OBSOLETE" )?`
  `( \s+ "SUP"` `\s+` **OIDs** ` )?`
  `( \s+ "EQUALITY"` `\s+` **OIDs** ` )?`
  `( \s+ "ORDERING"` `\s+` **OIDs** ` )?`
  `( \s+ "SUBSTR"` `\s+` **OIDs** ` )?`
  `( \s+ "SYNTAX"` `\s+` **NumericOIDwithLength** ` )?`
  `( \s+ "SINGLE-VALUE" )?`
  `( \s+ "COLLECTIVE" )?`
  `( \s+ "NO-USER-MODIFICATION" )?`
  `( \s+ "USAGE" \s+ ( "userApplications" | "directoryOperation" | "distributedOperation" | "dSAOperation" ) )?`
  `( \s+ ` **Extensions** ` )*`
  `\s* ")"`

## Matching Rules

* **MatchingRuleDescription** = `"(" \s*` **NumericOID**
  `( \s+ "NAME"` `\s+` **QuotedDescriptorS** ` )?`
  `( \s+ "DESC"` `\s+` **QuotedDString** ` )?`
  `( \s+ "OBSOLETE" )?`
  `\s+ "SYNTAX"` `\s+` **NumericOIDwithLength**
  `( \s+ ` **Extensions** ` )*`
  `\s* ")"`

## Matching Rule Uses

* **MatchingRuleUseDescription** = `"(" \s*` **NumericOID**
  `( \s+ "NAME"` `\s+` **QuotedDescriptorS** ` )?`
  `( \s+ "DESC"` `\s+` **QuotedDString** ` )?`
  `( \s+ "OBSOLETE" )?`
  `\s+ "APPLIES"` `\s+` **OIDs**
  `( \s+ ` **Extensions** ` )*`
  `\s* ")"`

## LDAP Syntax

* **SyntaxDescription** = `"(" \s*` **NumericOID**
  `( \s+ "DESC"` `\s+` **QuotedDString** ` )?`
  `( \s+ ` **Extensions** ` )*`
  `\s* ")"`

## DIT Content Rules

* **DITContentRuleDescription** = `"(" \s*` **NumericOID**
  `( \s+ "NAME"` `\s+` **QuotedDescriptorS** ` )?`
  `( \s+ "DESC"` `\s+` **QuotedDString** ` )?`
  `( \s+ "OBSOLETE" )?`
  `( \s+ "AUX"` `\s+` **OIDs** ` )?`
  `( \s+ "MUST"` `\s+` **OIDs** ` )?`
  `( \s+ "MAY"` `\s+` **OIDs** ` )?`
  `( \s+ "NOT"` `\s+` **OIDs** ` )?`
  `( \s+ ` **Extensions** ` )*`
  `\s* ")"`

## DIT Structure Rules

DIT structure rules does not starting with **NumericOID**.

* **RuleIDS** =
    - `[0-9]+` |
    - `"(" \s* [0-9]+ ( \s+ [0-9]+ )* \s* ")"`
* **DITStructureRuleDescription** = `"(" \s* [0-9]+`
  `( \s+ "NAME"` `\s+` **QuotedDescriptorS** ` )?`
  `( \s+ "DESC"` `\s+` **QuotedDString** ` )?`
  `( \s+ "OBSOLETE" )?`
  `\s+ "FORM"` `\s+` **OIDs**
  `( \s+ "SUP"` `\s+` **RuleIDS** ` )?`
  `( \s+ ` **Extensions** ` )*`
  `\s* ")"`

## Name Forms

* **NameFormDescription** = `"(" \s*` **NumericOID**
  `( \s+ "NAME"` `\s+` **QuotedDescriptorS** ` )?`
  `( \s+ "DESC"` `\s+` **QuotedDString** ` )?`
  `( \s+ "OBSOLETE" )?`
  `\s+ "FORM"` `\s+` **OID**
  `\s+ "MUST"` `\s+` **OIDs**
  `( \s+ "MAY"` `\s+` **OIDs** ` )?`
  `( \s+ ` **Extensions** ` )*`
  `\s* ")"`
