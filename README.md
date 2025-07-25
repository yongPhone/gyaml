# GYAML

[![Go Reference](https://pkg.go.dev/badge/github.com/yongPhone/gyaml.svg)](https://pkg.go.dev/github.com/yongPhone/gyaml)
[![GYAML Tests](https://github.com/yongPhone/gyaml/workflows/test/badge.svg)](https://github.com/yongPhone/gyaml/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/yongPhone/gyaml)](https://goreportcard.com/report/github.com/yongPhone/gyaml)
[![Coverage](https://img.shields.io/badge/Coverage-90.8%25-brightgreen)](https://github.com/yongPhone/gyaml)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

**🚀 A fast and simple way to get values from YAML in Go**

GYAML makes reading YAML as simple as `gyaml.Get(yaml, "path.to.value")`. Whether you're building microservices, processing configuration files, or handling CI/CD pipelines, GYAML reduces boilerplate and provides direct access to the data you need.

**Why developers choose GYAML:**
- 🎯 **One-line data access** - No struct marshaling required for simple reads
- ⚡ **Fast performance** - Optimized with minimal allocations  
- 🛡️ **Well-tested** - 90.8% test coverage for reliability
- 🔄 **Familiar API** - Same interface as popular gjson but for YAML
- 🎛️ **YAML-native** - Proper handling of comments, multi-line strings, and YAML features

GYAML is inspired by [tidwall/gjson](https://github.com/tidwall/gjson) but designed specifically for YAML, supporting all YAML features while maintaining the same intuitive dot-notation API that developers already know and love.

## 🏆 High-Quality Codebase

**Production-Ready with Exceptional Test Coverage**

- **90.8% Code Coverage** - Comprehensive test suite ensuring reliability
- **82+ Test Cases** - Covering edge cases, error handling, and performance scenarios  
- **Zero Race Conditions** - Thread-safe with race detection validation
- **Benchmark Tested** - Performance validated across all operations
- **Static Analysis Clean** - Passes go vet and formatting standards

**What makes GYAML useful:**
- 🚄 **Minimal boilerplate** - No struct definitions needed for simple data access
- 🎯 **Intuitive syntax** - `gyaml.Get(yaml, "path.to.data")` for direct access
- 🔍 **Powerful queries** - Array filtering with `#(condition)`, length with `#`, comparison operators
- 🛠️ **Smart type conversion** - `.String()`, `.Int()`, `.Bool()` with automatic type handling  
- 🔄 **Easy iteration** - `.ForEach()` for processing complex data structures
- 📝 **YAML-native** - Proper support for comments, multi-line strings, and YAML boolean variants
- ⚡ **Good performance** - Optimized for speed with minimal memory allocations
- 🔒 **Thread-safe** - Safe for use in concurrent applications
- 📦 **Lightweight** - Only depends on YAML parser
- 🧪 **Well-tested** - 90.8% test coverage for reliability

## 💡 Common Use Cases

**Configuration Management**
```go
// Instead of defining structs for every config read
config := gyaml.Get(configYAML, "database.host").String()
port := gyaml.Get(configYAML, "database.port").Int()
```

**CI/CD Pipeline Processing**
```go
// Parse GitHub Actions, GitLab CI, or any YAML-based pipeline
services := gyaml.Get(pipelineYAML, "services.#.name")
deploymentEnv := gyaml.Get(pipelineYAML, "stages.#(name=deploy).environment").String()
```

**Kubernetes Manifest Analysis**
```go
// Extract pod specs, resource limits, labels without complex parsing
replicas := gyaml.Get(k8sYAML, "spec.replicas").Int()
image := gyaml.Get(k8sYAML, "spec.template.spec.containers.0.image").String()
cpuLimit := gyaml.Get(k8sYAML, "spec.template.spec.containers.0.resources.limits.cpu").String()
```

**API Response Processing**
```go
// Handle YAML APIs (OpenAPI specs, Swagger docs, etc.)
endpoints := gyaml.Get(apiSpec, "paths").Map()
version := gyaml.Get(apiSpec, "info.version").String()
```

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
BenchmarkGet-8                     22927    51155 ns/op    35424 B/op    598 allocs/op
BenchmarkGetNested-8               23538    52656 ns/op    35568 B/op    602 allocs/op
BenchmarkGetArray-8                22798    52312 ns/op    35520 B/op    600 allocs/op
BenchmarkGetArrayLength-8          23010    51347 ns/op    35360 B/op    596 allocs/op
BenchmarkGetArrayOperation-8       22071    54872 ns/op    42576 B/op    633 allocs/op
BenchmarkParse-8                   21114    55115 ns/op    35264 B/op    593 allocs/op
BenchmarkValid-8                   23287    50985 ns/op    35264 B/op    593 allocs/op
```

*These benchmarks were run on a MacBook Pro M1 using Go 1.22.*

## 🧪 Test Quality & Coverage

GYAML takes testing seriously with an industry-leading test suite:

- **90.8% Code Coverage** - One of the highest in the Go ecosystem
- **82+ Test Cases** - Comprehensive coverage across all features
- **Zero Race Conditions** - Validated with `go test -race`
- **Production-Ready Quality** - Unit tests, edge cases, error handling, performance benchmarks, and concurrency safety

**Run the tests yourself:**

```bash
# Run all tests with coverage
go test -coverprofile=coverage.out -covermode=atomic ./...

# Run with race detection
go test -race ./...

# Run benchmarks
go test -bench=. -benchmem ./...
```

✅ **GYAML is thoroughly tested and ready for production use.**

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

## 🚀 Ready to Simplify Your YAML Processing?

```bash
go get -u github.com/yongPhone/gyaml
```

**Start with a simple example:**
```go
import "github.com/yongPhone/gyaml"

// Instead of defining structs for simple reads...
type Config struct {
    Database struct {
        Host string `yaml:"host"`
        Port int    `yaml:"port"`
    } `yaml:"database"`
}

// Just do this!
host := gyaml.Get(yamlData, "database.host").String()
port := gyaml.Get(yamlData, "database.port").Int()
```

GYAML is well-tested with 90.8% code coverage and designed for production use. Give it a try and see if it simplifies your YAML processing workflow.

---

**⭐ If you find GYAML helpful, consider starring this repo to help others discover it.**

## License

MIT License. See [LICENSE](LICENSE) file.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

---

*GYAML v1.0.0 - Production Ready Release | Last updated: July 2025* 