# 4GL syntax
The 4GL is not case-sensitive, although lowercase is preferable.
For the annotator several language constructs are relevant:

## class statement
```
CLASS class-type-name [ INHERITS super-type-name]
  [ IMPLEMENTS interface-type-name [ , interface-type-name ] ... ]
  [ USE-WIDGET-POOL ]
  [ ABSTRACT | FINAL ]
  [ SERIALIZABLE ]:
  class-body
```

classes are ended with either `end.` or `end class.` 
A class found in `misc/string/Hello.cls` has a class type name of `misc.string.Hello`.
File names (and therefor class names) *are* case-sensitive.


## method statement
Methods in classes are eexpressed like:
```
METHOD [ PRIVATE | PACKAGE-PRIVATE | PROTECTED | PACKAGE-PROTECTED | PUBLIC ]
  [ STATIC | ABSTRACT ] [ OVERRIDE ][ FINAL ] 
  { VOID |return-type}method-name 
  ( [ parameter [ , parameter ]...] ) :
```

methods are ended with either `end.` or `end method.`

## annotations
```
@annotation[(attribute = "value"[,attribute = "value"]...)].
```
Where `annotation` is the name of the annotation. An example would be:
```
@todo(version="13.0", what="check generic constructs").
```
Annotations can span multiple line. It ends when the closing period is found.

annotation: The annotation's name can be any character string that you choose.

attribute: The attribute's name can be any character string that you choose. Attribute/value pairs are optional. Attributes are not case-sensitive.

value: The value can be any character string that you choose. Attribute/value pairs are optional. 

The escape character in the 4GL is the tilde. Escaped chars are 
- `~`
- `"`
- `{`
- `~n`
  

## comments
the 4GL uses c++ style comments, so either `/* comment */` or `// comment`. Comments can be nested.

## output
Assume a class file `misc/string/Hello.cls`:
```
@todo(version="13.0", what="refector generics").
class misc.string.Hello:

  @http(method="get", comment="blabla").
  method public void GetHellos():

  end.

end class.
```

has an output of:
```
{
  "annotations": {
    "todo": [
      {}
    ],
    "http": [
      {}
    ]
  }
}
```

where the annotation looks like:
```
{
  "name": "todo",
  "attributes": [
    { "name": "version", "value": "13.0"},
    { "name": "what", "value": "refector generics"},
  ],
  "file": "misc/string/Hello.cls",
  "classname": "misc.string.Hello",
  "type": "class",
  "annotationLine": 1,
  "constructLine": 2
}
```
or 
```
{
  "name": "http",
  "attributes": [
    { "name": "method", "value": "get"},
    { "name": "comment", "value": "blabla"},
  ],
  "file": "misc/string/Hello.cls",
  "classname": "misc.string.Hello",
  "type": "method",
  "constructName": "GetHellos",
  "annotationLine": 4,
  "constructLine": 5
}
``` 
