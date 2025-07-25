# GYAML Path Syntax

This document describes the complete path syntax for GYAML. A path is a series of keys separated by a dot. The value for a key is retrieved by using one of the "get" functions such as `Get`, `GetBytes`, or `Parse` followed by a `Get`.

## Basic Paths

A path is a series of keys separated by a dot. A key may contain special characters.

```go
"name"               >> "Janet"
"name.first"         >> "Janet"
"name.last"          >> "Prichard"
"age"                >> 47
```

## Arrays

Arrays are accessed by index or through various operations.

### Array indexing

Use the index as the key to access an array element:

```go
"children.0"         >> "Sara"
"children.1"         >> "Alex"
"children.2"         >> "Jack"
```

### Array length

Use the `#` character to get the number of elements in an array:

```go
"children.#"         >> 3
```

## Nested Arrays

You can access nested arrays and get values from all elements.

### Getting all values of a specific key

To get all values of a key from array elements:

```go
"children.#.name"    >> ["Sara","Alex","Jack"]
```

This returns an array containing the `name` value from each element in the `children` array.

### Nested object access

For arrays containing objects:

```yaml
friends:
  - first: "Dale"
    last: "Murphy"
    age: 44
  - first: "Roger" 
    last: "Craig"
    age: 68
```

```go
"friends.#.first"    >> ["Dale","Roger"]
"friends.0.age"      >> 44
"friends.1.last"     >> "Craig"
```

## Queries

You can query an array for objects that match specific conditions.

### Basic query syntax

The basic syntax for a query is `#(key=value)`:

```go
friends.#(last="Murphy").first   >> "Dale"
```

This finds the first object in the `friends` array where `last` equals "Murphy" and returns the `first` value.

### Comparison operators

Currently supported operators:

- `=` - Equals
- `!=` - Not equals  
- `<` - Less than
- `<=` - Less than or equal
- `>` - Greater than
- `>=` - Greater than or equal

Examples:

```go
friends.#(age>40).last        >> "Murphy"
friends.#(age>=68).first      >> "Roger"
```

### String matching

String values should be quoted:

```go
friends.#(first="Dale").age   >> 44
```

## YAML-Specific Features

GYAML supports YAML-specific syntax that differs from JSON.

### Multi-line strings

YAML supports multi-line strings using `|` (literal) and `>` (folded):

```yaml
description: |
  This is a multi-line string
  that preserves line breaks
summary: >
  This is a folded string
  that becomes one line
```

```go
"description"        >> "This is a multi-line string\nthat preserves line breaks"
"summary"           >> "This is a folded string that becomes one line"
```

### Boolean values

YAML supports various boolean representations:

```yaml
flags:
  debug: true
  verbose: false
  enabled: yes
  disabled: no
  active: on
  inactive: off
```

All of these can be converted to boolean values:

```go
"flags.debug"        >> true
"flags.enabled"      >> true (converted from "yes")
"flags.active"       >> true (converted from "on")
```

### Comments

Comments are ignored during parsing:

```yaml
# This is a comment
name: "John"  # Inline comment
age: 30
```

```go
"name"               >> "John"
"age"                >> 30
```

## Special Characters in Keys

Keys containing special characters or spaces should be handled carefully:

```yaml
"key-with-dashes": "value1"
"key_with_underscores": "value2"
"123numeric_key": "value3"
"key with spaces": "value4"
```

These can be accessed normally:

```go
"key-with-dashes"           >> "value1"
"key_with_underscores"      >> "value2"
"123numeric_key"            >> "value3"
"key with spaces"           >> "value4"
```

## Escaping

Special characters in values are automatically handled by the YAML parser:

```yaml
text: "Line 1\nLine 2\tTabbed"
path: "C:\\Windows\\System32"
unicode: "Hello \u4E16\u754C"
```

```go
"text"               >> "Line 1\nLine 2\tTabbed"
"path"               >> "C:\\Windows\\System32"
"unicode"            >> "Hello 世界"
```

## Error Handling

When a path doesn't exist, GYAML returns a `Result` with type `Null`:

```go
result := gyaml.Get(yaml, "nonexistent.path")
result.Exists()      >> false
result.Type          >> gyaml.Null
```

## Performance Considerations

- Paths are parsed once per `Get` call
- Deeper paths require more processing
- Array operations (`#`) are generally fast
- Complex queries may be slower on large arrays

## Examples

Here are some complete examples showing various path syntax:

```yaml
app:
  name: "My App"
  version: "1.0.0"
  features:
    - name: "feature1"
      enabled: true
    - name: "feature2"
      enabled: false
servers:
  - host: "server1.com"
    port: 8080
    roles: ["web", "api"]
  - host: "server2.com" 
    port: 8081
    roles: ["worker"]
```

```go
// Basic access
gyaml.Get(yaml, "app.name")                    // "My App"
gyaml.Get(yaml, "app.version")                 // "1.0.0"

// Array access
gyaml.Get(yaml, "servers.0.host")              // "server1.com"
gyaml.Get(yaml, "servers.#")                   // 2

// Array operations
gyaml.Get(yaml, "servers.#.host")              // ["server1.com", "server2.com"]
gyaml.Get(yaml, "app.features.#.name")         // ["feature1", "feature2"]

// Queries
gyaml.Get(yaml, `servers.#(port=8080).host`)   // "server1.com"
gyaml.Get(yaml, `app.features.#(enabled=true).name`) // "feature1"

// Nested access
gyaml.Get(yaml, "servers.0.roles.1")           // "api"
``` 