Here is syntax of LDAP schema in regular expression for understanding.

# Common Production

* **NumericOID** = `[0-9]+` `( \.[0-9]+ )*`
* **Descriptor** = `[a-zA-Z]` `[a-zA-Z0-9-]*`
* **OID** = **NumericOID** | **Descriptor**

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

## Object Class

* **ObjectClassDescription** = `"(" \s*` **NumericOID**
  `( \s+` `"NAME"` `\s+` **QuotedDescriptorS** ` )?`
  `( \s+` `"DESC"` `\s+` **QuotedDString** ` )?`
  `( \s+` `"OBSOLETE"` `\s+` **QuotedDString** ` )?`
  `( \s+` `"SUP"` `\s+` **OIDs** ` )?`
  `( \s+ ( "ABSTRACT" | "STRUCTURAL" | "AUXILIARY" ) )?`
  `( \s+` `"MUST"` `\s+` **OIDs** ` )?`
  `( \s+` `"MAY"` `\s+` **OIDs** ` )?`
  `( \s+ "X-" [a-zA-Z-_]+ \s+ ` **QuotedDStringS** ` )*`
  `\s* ")"`
