// Package gyaml provides a fast and simple way to get values from a YAML document.
// It has features such as one line retrieval, dot notation paths, iteration, and parsing YAML documents.
//
// This package is inspired by tidwall/gjson but works with YAML instead of JSON.
// GYAML supports YAML-specific features like multi-line strings, comments, and various boolean representations.
package gyaml

import (
	"fmt"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

// Type represents a YAML value type
type Type int

const (
	// Null represents a null YAML value
	Null Type = iota
	// False represents a false YAML value
	False
	// Number represents a numeric YAML value
	Number
	// String represents a string YAML value
	String
	// True represents a true YAML value
	True
	// YAML represents a YAML object
	YAML
)

// Result represents a YAML value that is returned from Get()
type Result struct {
	// Type is the YAML type
	Type Type
	// Raw is the raw YAML value
	Raw string
	// Str is the YAML string
	Str string
	// Num is the YAML number
	Num float64
	// Index of raw value in original YAML, or -1
	Index int
}

// String returns a string representation of the value.
func (t Result) String() string {
	switch t.Type {
	default:
		return ""
	case String:
		return t.Str
	case Number:
		if len(t.Raw) == 0 {
			return strconv.FormatFloat(t.Num, 'f', -1, 64)
		}
		return t.Raw
	case YAML:
		return t.Raw
	case True:
		return "true"
	case False:
		return "false"
	}
}

// Bool returns a boolean representation of the value.
func (t Result) Bool() bool {
	switch t.Type {
	default:
		return false
	case True:
		return true
	case String:
		lower := strings.ToLower(t.Str)
		// Handle YAML boolean-like strings
		switch lower {
		case "true", "yes", "on", "1":
			return true
		case "false", "no", "off", "0", "":
			return false
		default:
			// Try standard parsing
			b, _ := strconv.ParseBool(lower)
			return b
		}
	case Number:
		return t.Num != 0
	}
}

// Int returns an integer representation of the value.
func (t Result) Int() int64 {
	switch t.Type {
	default:
		return 0
	case True:
		return 1
	case String:
		n, _ := strconv.ParseInt(t.Str, 10, 64)
		return n
	case Number:
		// Check if we can parse from Raw to avoid float64 precision loss
		if t.Raw != "" {
			if n, err := strconv.ParseInt(strings.TrimSpace(t.Raw), 10, 64); err == nil {
				return n
			}
		}
		return int64(t.Num)
	}
}

// Uint returns an unsigned integer representation of the value.
func (t Result) Uint() uint64 {
	switch t.Type {
	default:
		return 0
	case True:
		return 1
	case String:
		n, _ := strconv.ParseUint(t.Str, 10, 64)
		return n
	case Number:
		return uint64(t.Num)
	}
}

// Float returns a float64 representation of the value.
func (t Result) Float() float64 {
	switch t.Type {
	default:
		return 0
	case True:
		return 1
	case String:
		n, _ := strconv.ParseFloat(t.Str, 64)
		return n
	case Number:
		return t.Num
	}
}

// Array returns an array of values.
func (t Result) Array() []Result {
	if t.Type != YAML {
		return nil
	}
	var any interface{}
	if err := yaml.Unmarshal([]byte(t.Raw), &any); err != nil {
		return nil
	}
	arr, ok := any.([]interface{})
	if !ok {
		return nil
	}
	results := make([]Result, len(arr))
	for i, v := range arr {
		results[i] = makeResult(v)
	}
	return results
}

// Map returns a map of key-value pairs.
func (t Result) Map() map[string]Result {
	if t.Type != YAML {
		return nil
	}
	var any interface{}
	if err := yaml.Unmarshal([]byte(t.Raw), &any); err != nil {
		return nil
	}
	obj, ok := any.(map[string]interface{})
	if !ok {
		return nil
	}
	results := make(map[string]Result)
	for k, v := range obj {
		results[k] = makeResult(v)
	}
	return results
}

// Get returns the result for the specified path.
func (t Result) Get(path string) Result {
	if t.Type != YAML {
		return Result{}
	}
	return Get(t.Raw, path)
}

// Value returns the raw interface{} value.
func (t Result) Value() interface{} {
	if t.Type == YAML {
		var any interface{}
		yaml.Unmarshal([]byte(t.Raw), &any)
		return any
	}
	switch t.Type {
	default:
		return nil
	case False:
		return false
	case Number:
		return t.Num
	case String:
		return t.Str
	case True:
		return true
	}
}

// Exists returns true if value exists.
func (t Result) Exists() bool {
	return t.Type != Null
}

// ForEach iterates through values.
func (t Result) ForEach(iterator func(key, value Result) bool) {
	if !t.Exists() {
		return
	}
	if t.Type != YAML {
		return
	}
	var any interface{}
	if err := yaml.Unmarshal([]byte(t.Raw), &any); err != nil {
		return
	}
	switch obj := any.(type) {
	case map[string]interface{}:
		for k, v := range obj {
			if !iterator(Result{Type: String, Str: k}, makeResult(v)) {
				return
			}
		}
	case []interface{}:
		for i, v := range obj {
			if !iterator(Result{Type: Number, Num: float64(i)}, makeResult(v)) {
				return
			}
		}
	}
}

// makeResult creates a Result from an interface{} value
func makeResult(value interface{}) Result {
	if value == nil {
		return Result{Type: Null}
	}

	switch v := value.(type) {
	case bool:
		if v {
			return Result{Type: True}
		}
		return Result{Type: False}
	case string:
		return Result{Type: String, Str: v}
	case int:
		return Result{Type: Number, Num: float64(v), Raw: strconv.FormatInt(int64(v), 10)}
	case int8:
		return Result{Type: Number, Num: float64(v), Raw: strconv.FormatInt(int64(v), 10)}
	case int16:
		return Result{Type: Number, Num: float64(v), Raw: strconv.FormatInt(int64(v), 10)}
	case int32:
		return Result{Type: Number, Num: float64(v), Raw: strconv.FormatInt(int64(v), 10)}
	case int64:
		return Result{Type: Number, Num: float64(v), Raw: strconv.FormatInt(v, 10)}
	case uint:
		return Result{Type: Number, Num: float64(v), Raw: strconv.FormatUint(uint64(v), 10)}
	case uint8:
		return Result{Type: Number, Num: float64(v), Raw: strconv.FormatUint(uint64(v), 10)}
	case uint16:
		return Result{Type: Number, Num: float64(v), Raw: strconv.FormatUint(uint64(v), 10)}
	case uint32:
		return Result{Type: Number, Num: float64(v), Raw: strconv.FormatUint(uint64(v), 10)}
	case uint64:
		return Result{Type: Number, Num: float64(v), Raw: strconv.FormatUint(v, 10)}
	case float32:
		return Result{Type: Number, Num: float64(v), Raw: strconv.FormatFloat(float64(v), 'g', -1, 32)}
	case float64:
		return Result{Type: Number, Num: v, Raw: strconv.FormatFloat(v, 'g', -1, 64)}
	default:
		// For complex types, marshal back to YAML
		raw, err := yaml.Marshal(v)
		if err != nil {
			return Result{Type: Null}
		}
		return Result{Type: YAML, Raw: string(raw)}
	}
}

// Get searches YAML for the specified path.
// A path is in dot syntax, such as "name.last" or "age".
// When the value is found it's returned immediately.
func Get(yamlStr, path string) Result {
	if len(yamlStr) == 0 {
		return Result{Type: Null}
	}

	var root interface{}
	if err := yaml.Unmarshal([]byte(yamlStr), &root); err != nil {
		return Result{Type: Null}
	}

	// If path is empty, return the entire document
	if len(path) == 0 {
		return Result{Type: YAML, Raw: yamlStr}
	}

	return getByPath(root, path)
}

// GetBytes searches YAML bytes for the specified path.
func GetBytes(yamlBytes []byte, path string) Result {
	return Get(string(yamlBytes), path)
}

// Parse parses the YAML and returns a result.
func Parse(yamlStr string) Result {
	if len(yamlStr) == 0 {
		return Result{Type: Null}
	}

	var root interface{}
	if err := yaml.Unmarshal([]byte(yamlStr), &root); err != nil {
		return Result{Type: Null}
	}

	return Result{Type: YAML, Raw: yamlStr}
}

// Valid returns true if the YAML is valid.
func Valid(yamlStr string) bool {
	var root interface{}
	return yaml.Unmarshal([]byte(yamlStr), &root) == nil
}

// getByPath navigates through the parsed YAML structure using the path
func getByPath(root interface{}, path string) Result {
	if path == "" {
		// For empty path, return the root as YAML type if it's complex
		if root == nil {
			return Result{Type: Null}
		}
		switch v := root.(type) {
		case string, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64, bool:
			return makeResult(v)
		default:
			// Complex type, marshal to YAML and return as YAML type
			raw, err := yaml.Marshal(v)
			if err != nil {
				return Result{Type: Null}
			}
			return Result{Type: YAML, Raw: string(raw)}
		}
	}

	parts := strings.Split(path, ".")
	current := root

	for i, part := range parts {
		if part == "" {
			continue
		}

		// Handle array length with #
		if part == "#" {
			// Check if this is the last part or if next part is empty
			if i == len(parts)-1 {
				switch v := current.(type) {
				case []interface{}:
					return Result{Type: Number, Num: float64(len(v))}
				case map[string]interface{}:
					return Result{Type: Number, Num: float64(len(v))}
				default:
					return Result{Type: Null}
				}
			} else {
				// This is #.something, collect remaining path and handle array operation
				remainingPath := strings.Join(parts[i+1:], ".")
				return handleArrayOperation(current, remainingPath)
			}
		}

		// Handle array queries like #(key=value)
		if strings.HasPrefix(part, "#(") && strings.HasSuffix(part, ")") {
			query := part[2 : len(part)-1] // Remove #( and )
			return handleArrayQuery(current, query)
		}

		// Handle array access with wildcard or specific operations that start with #
		if strings.HasPrefix(part, "#") && part != "#" {
			// Handle cases like #key
			remaining := part[1:]
			return handleArrayOperation(current, remaining)
		}

		// Handle array index
		if idx, err := strconv.Atoi(part); err == nil {
			switch v := current.(type) {
			case []interface{}:
				if idx < 0 || idx >= len(v) {
					return Result{Type: Null}
				}
				current = v[idx]
				continue
			default:
				return Result{Type: Null}
			}
		}

		// Handle map access
		switch v := current.(type) {
		case map[string]interface{}:
			val, exists := v[part]
			if !exists {
				return Result{Type: Null}
			}
			current = val
		case map[interface{}]interface{}:
			val, exists := v[part]
			if !exists {
				return Result{Type: Null}
			}
			current = val
		default:
			return Result{Type: Null}
		}
	}

	return makeResult(current)
}

// handleArrayQuery handles queries like #(key=value)
func handleArrayQuery(current interface{}, query string) Result {
	arr, ok := current.([]interface{})
	if !ok {
		return Result{Type: Null}
	}

	// Parse the query (simple key=value for now)
	parts := strings.SplitN(query, "=", 2)
	if len(parts) != 2 {
		return Result{Type: Null}
	}

	key := parts[0]
	value := strings.Trim(parts[1], `"'`)

	for _, item := range arr {
		if obj, ok := item.(map[string]interface{}); ok {
			if val, exists := obj[key]; exists {
				if fmt.Sprintf("%v", val) == value {
					return makeResult(item)
				}
			}
		}
	}

	return Result{Type: Null}
}

// handleArrayOperation handles operations like #.key (get all values of key from array elements)
func handleArrayOperation(current interface{}, path string) Result {
	arr, ok := current.([]interface{})
	if !ok {
		return Result{Type: Null}
	}

	if path == "" {
		// Return the whole array
		return makeResult(arr)
	}

	var results []interface{}
	for _, item := range arr {
		// For each item in the array, get the value at the specified path
		itemResult := getByPath(item, path)
		if itemResult.Exists() {
			results = append(results, itemResult.Value())
		}
	}

	return makeResult(results)
}

// ForEachLine iterates through each line of a YAML document.
func ForEachLine(yamlStr string, iterator func(line Result) bool) {
	lines := strings.Split(yamlStr, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if !iterator(Parse(line)) {
			return
		}
	}
}
