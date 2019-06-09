# Build

The following command build utilities:

```sh
go build github.com/yinyin/go-ldap-schema-parser/cmd/ldap-schema-parser
go build github.com/yinyin/go-ldap-schema-parser/cmd/ldif-subschema-extract
go build github.com/yinyin/go-ldap-schema-parser/cmd/rfc-ldap-schema-extract
```

# Import Schema Elements

```sh
./rfc-ldap-schema-extract -out /tmp/ldap-schema-elements.txt \
    -rfc4512 docs/spec/rfc4512.txt \
    -rfc4517 docs/spec/rfc4517.txt \
    -rfc4519 docs/spec/rfc4519.txt \
    -rfc4523 docs/spec/rfc4523.txt
```

```sh
./ldif-subschema-extract -out /tmp/ldap-schema-elements.txt \
    docs/schema/core.ldif \
    docs/schema/cosine.ldif \
    docs/schema/inetorgperson.ldif \
    docs/schema/nis.ldif
```
