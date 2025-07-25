# GYAML

<a href="https://pkg.go.dev/github.com/yongPhone/gyaml"><img src="https://pkg.go.dev/badge/github.com/yongPhone/gyaml.svg" alt="Go Reference"></a>
<a href="https://github.com/yongPhone/gyaml/actions"><img src="https://github.com/yongPhone/gyaml/workflows/test/badge.svg" alt="GYAML Tests"></a>
<a href="https://goreportcard.com/report/github.com/yongPhone/gyaml"><img src="https://goreportcard.com/badge/github.com/yongPhone/gyaml" alt="Go Report Card"></a>

**get yaml values quickly**

GYAML is a Go package that provides a fast and simple way to get values from a YAML document. It has features such as one line retrieval, dot notation paths, iteration, and parsing YAML documents.

GYAML is inspired by [tidwall/gjson](https://github.com/tidwall/gjson) but works with YAML instead of JSON, providing the same intuitive API while supporting YAML-specific features like comments, multi-line strings, and various boolean representations.

**Features:**
- Fast and simple value retrieval
- Dot notation path syntax
- Array operations and queries
- Type conversion methods
- Iteration support
- YAML-specific features (comments, multi-line strings, boolean variants)
- Thread-safe operations
- Zero external dependencies except YAML parser

## Getting Started

### Installing

To start using GYAML, install Go and run `go get`:

```bash
$ go get -u github.com/yongPhone/gyaml
```

This will retrieve the library.

### Get a value

Get searches YAML for the specified path. A path is in dot syntax, such as "name.last" or "age". When the value is found it's returned immediately.

```go
package main

import "github.com/yongPhone/gyaml"

const yaml = `
name:
  first: "Janet"
  last: "Prichard"
age: 47
`

func main() {
    value := gyaml.Get(yaml, "name.last")
    println(value.String())
}
```

This will print:
```
Prichard
```

*There's also GetBytes for working with YAML byte slices.*

## Path Syntax

Below is a quick overview of the path syntax. A path is a series of keys separated by a dot. A key may contain special characters. To access an array value use the index as the key. To get the number of elements in an array or to access a child path, use the '#' character.

```yaml
name:
  first: "Tom"
  last: "Anderson"
age: 37
children:
  - "Sara"
  - "Alex" 
  - "Jack"
favorite_movie: "Deer Hunter"
friends:
  - first: "Dale"
    last: "Murphy"
    age: 44
    networks: ["ig", "fb", "tw"]
  - first: "Roger"
    last: "Craig"
    age: 68
    networks: ["fb", "tw"]
  - first: "Jane"
    last: "Murphy"
    age: 47
    networks: ["ig", "tw"]
```

```go
"name.last"          >> "Anderson"
"age"                >> 37
"children"           >> ["Sara","Alex","Jack"]
"children.#"         >> 3
"children.1"         >> "Alex"
"children.#.length"  >> [4,4,4]
"friends.#.first"    >> ["Dale","Roger","Jane"]
```

### Arrays

Arrays are accessed by index or through various operations:

```go
"children.0"           >> "Sara"
"children.1"           >> "Alex"
"children.#"           >> 3
"friends.#.age"        >> [44,68,47]
```

### Queries

You can also query an object inside an array:

```go
friends.#(last="Murphy").first   >> "Dale"
friends.#(age>65).last           >> "Craig"
```

## Result Type

All `Get` methods return a `Result` type. The `Result` type has several methods:

```go
result.Type      // Returns the YAML type (Null, False, Number, String, True, YAML)
result.Exists()  // Returns true if the value exists
result.String()  // Returns a string representation
result.Int()     // Returns an int64 representation
result.Uint()    // Returns a uint64 representation  
result.Float()   // Returns a float64 representation
result.Bool()    // Returns a bool representation
result.Array()   // Returns an array of Result values
result.Map()     // Returns a map[string]Result
result.Value()   // Returns the raw interface{} value
result.Raw       // Returns the raw YAML value as a string
```

### 64-bit integers

The `Int()` and `Uint()` methods return 64-bit integers:

```go
result.Int()   // int64
result.Uint()  // uint64
```

### Boolean Values

GYAML supports various boolean representations common in YAML:

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
gyaml.Get(yaml, "flags.debug").Bool()    // true
gyaml.Get(yaml, "flags.enabled").Bool()  // true
gyaml.Get(yaml, "flags.active").Bool()   // true
```

## YAML-Specific Features

### Multi-line Strings

GYAML properly handles YAML's multi-line string syntax:

```yaml
description: |
  This is a multi-line string
  that preserves line breaks
  and supports formatting
```

### Comments

Comments in YAML are preserved during parsing and don't affect value retrieval:

```yaml
# Database configuration
database:
  host: "localhost"  # Local development
  port: 5432
```

## Get nested array values

Suppose you want all the last names from the following YAML:

```yaml
programmers:
  - firstName: "Janet"
    lastName: "McLaughlin"
  - firstName: "Elliotte" 
    lastName: "Hunter"
  - firstName: "Jason"
    lastName: "Harold"
```

You would use the path "programmers.#.lastName" like such:

```go
result := gyaml.Get(yaml, "programmers.#.lastName")
for _, name := range result.Array() {
    println(name.String())
}
```

You can also query an object inside an array:

```go
name := gyaml.Get(yaml, `programmers.#(lastName="Hunter").firstName`)
println(name.String())  // prints "Elliotte"
```

## Iterate through an object or array

The `ForEach` function allows for quickly iterating through an object or array. The key and value are passed to the iterator function for objects. Only the value is passed for arrays. Returning `false` from an iterator will stop iteration.

```go
result := gyaml.Get(yaml, "programmers")
result.ForEach(func(key, value gyaml.Result) bool {
    println(value.String()) 
    return true // keep iterating
})
```

## Simple Parse and Get

There's a `Parse(yaml)` function that will do a simple parse, and `result.Get(path)` that will search a result.

For example, all of these will return the same result:

```go
gyaml.Parse(yaml).Get("name").Get("last")
gyaml.Get(yaml, "name").Get("last")
gyaml.Get(yaml, "name.last")
```

## Check for the existence of a value

Sometimes you just want to know if a value exists:

```go
value := gyaml.Get(yaml, "name.last")
if !value.Exists() {
    println("no last name")
} else {
    println(value.String())
}

// Or as one step
if gyaml.Get(yaml, "name.last").Exists() {
    println("has a last name")
}
```

## Validate YAML

The `Get*` and `Parse*` functions expect that the YAML is well-formed. Bad YAML will not panic, but it may return back unexpected results.

If you are consuming YAML from an unpredictable source then you may want to validate prior to using GYAML:

```go
if !gyaml.Valid(yaml) {
    return errors.New("invalid yaml")
}
value := gyaml.Get(yaml, "name.last")
```

## Working with Bytes

If your YAML is contained in a `[]byte` slice, there's the GetBytes function. This is preferred over `Get(string(data), path)`:

```go
var yaml []byte = ...
result := gyaml.GetBytes(yaml, path)
```

## Performance

GYAML is designed for performance. Here are some benchmark results:

```
BenchmarkGet-8                     22401    54636 ns/op    35424 B/op    598 allocs/op
BenchmarkGetNested-8               21729    52406 ns/op    35568 B/op    602 allocs/op
BenchmarkGetArray-8                21033    52308 ns/op    35520 B/op    600 allocs/op
BenchmarkGetArrayLength-8          22802    64219 ns/op    35360 B/op    596 allocs/op
BenchmarkGetArrayOperation-8       21384    65923 ns/op    42576 B/op    633 allocs/op
BenchmarkParse-8                   23568    52108 ns/op    35264 B/op    593 allocs/op
BenchmarkValid-8                   22434    51725 ns/op    35264 B/op    593 allocs/op
```

*These benchmarks were run on a MacBook Pro M1 using Go 1.19.*

## API Compatibility

GYAML provides the same API as gjson, making it easy to switch between JSON and YAML:

| Feature | gjson (JSON) | gyaml (YAML) |
|---------|-------------|-------------|
| Basic path queries | ✅ | ✅ |
| Array operations | ✅ | ✅ |
| Conditional queries | ✅ | ✅ |  
| Type conversion | ✅ | ✅ |
| ForEach iteration | ✅ | ✅ |
| Multi-line strings | ❌ | ✅ |
| Comments support | ❌ | ✅ |
| Boolean variants | Partial | ✅ |

## License

MIT License. See [LICENSE](LICENSE) file.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request. 