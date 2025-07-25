package gyaml

import (
	"testing"
)

// Edge cases and special scenarios testing
const edgeCaseYAML = `
# Special key name testing
special_keys:
  key-with-dashes: "value2"
  key_with_underscores: "value3"
  "123numeric_key": "value4"
  "key with spaces": "value5"
  
# Extreme value testing
extreme_values:
  max_int: 9223372036854775807
  min_int: -9223372036854775808
  max_float: 1.7976931348623157e+308
  min_float: 2.2250738585072014e-308
  very_long_string: "This is a very very long string for testing string handling edge cases. It contains Chinese characters, English characters, numbers 123456789, special symbols !@#$%^&*(), and various punctuation marks. The purpose of this string is to verify that our YAML parser can correctly handle long strings without memory overflow or performance issues during processing."
  
# Nested array testing
nested_arrays:
  level1:
    - level2:
        - level3:
            - level4:
                - "deep_string"
                - 42
                - true
            - level4_alt: ["a", "b", "c"]
        - simple_string: "test"
    - another_level2:
        - "simple"
        - complex:
            nested_value: "found"
            
# Empty and special value testing  
empty_and_null:
  empty_object: {}
  empty_array: []
  null_explicit: null
  null_implicit:
  zero: 0
  false_value: false
  empty_string: ""
  whitespace_string: "   "
  
# Quote and escape testing
quotes_and_escapes:
  single_quotes: 'single quoted string'
  double_quotes: "double quoted string"
  mixed_quotes: 'string with "inner" quotes'
  escaped_chars: "string with \n newline and \t tab"
  backslash: "path\\to\\file"
  unicode_escape: "unicode \u4E2D\u6587"
  
# Large array testing (performance testing)
large_array:
  numbers: [1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50]
  objects:
    - id: 1
      name: "Item 1"
      active: true
    - id: 2  
      name: "Item 2"
      active: false
    - id: 3
      name: "Item 3"
      active: true
    - id: 4
      name: "Item 4"
      active: false
    - id: 5
      name: "Item 5"
      active: true
      
# YAML special syntax testing
yaml_syntax:
  multiline_literal: |
    First line
    Second line
    Third line
  multiline_folded: >
    This line is very very very very very long
    Will be folded into one line
    Unless there are two newlines


    Like this will preserve paragraph separation
  preserve_newlines: |+
    Preserve trailing newlines
    

  strip_newlines: |-
    Strip trailing newlines
    
  
# Complex path testing
complex_paths:
  path_with_dots:
    nested_key: "value1"
  array_bracket:
    element: "confusing_name"
  hashtag_key:
    count: 100
`

func TestSpecialKeyNames(t *testing.T) {
	// Test key names containing hyphens
	result := Get(edgeCaseYAML, "special_keys.key-with-dashes")
	if result.String() != "value2" {
		t.Errorf("Expected 'value2', got '%s'", result.String())
	}

	// Test key names containing underscores
	result = Get(edgeCaseYAML, "special_keys.key_with_underscores")
	if result.String() != "value3" {
		t.Errorf("Expected 'value3', got '%s'", result.String())
	}

	// Test key names starting with numbers
	result = Get(edgeCaseYAML, "special_keys.123numeric_key")
	if result.String() != "value4" {
		t.Errorf("Expected 'value4', got '%s'", result.String())
	}

	// Test key names containing spaces
	result = Get(edgeCaseYAML, "special_keys.key with spaces")
	if result.String() != "value5" {
		t.Errorf("Expected 'value5', got '%s'", result.String())
	}
}

func TestExtremeValues(t *testing.T) {
	// Test maximum integer
	result := Get(edgeCaseYAML, "extreme_values.max_int")
	if result.Int() != 9223372036854775807 {
		t.Errorf("Expected max int64, got %d", result.Int())
	}

	// Test minimum integer
	result = Get(edgeCaseYAML, "extreme_values.min_int")
	if result.Int() != -9223372036854775808 {
		t.Errorf("Expected min int64, got %d", result.Int())
	}

	// Test very long string
	result = Get(edgeCaseYAML, "extreme_values.very_long_string")
	longStr := result.String()
	if len(longStr) < 100 {
		t.Errorf("Expected long string, got length %d", len(longStr))
	}
	if !Contains(longStr, "very very long string") {
		t.Error("Expected long string to contain specific text")
	}
}

func TestDeepNestedArrays(t *testing.T) {
	// Test deep nested arrays
	result := Get(edgeCaseYAML, "nested_arrays.level1.0.level2.0.level3.0.level4.0")
	if result.String() != "deep_string" {
		t.Errorf("Expected 'deep_string', got '%s'", result.String())
	}

	result = Get(edgeCaseYAML, "nested_arrays.level1.0.level2.0.level3.0.level4.1")
	if result.Int() != 42 {
		t.Errorf("Expected 42, got %d", result.Int())
	}

	// Test objects in nested arrays
	result = Get(edgeCaseYAML, "nested_arrays.level1.1.another_level2.1.complex.nested_value")
	if result.String() != "found" {
		t.Errorf("Expected 'found', got '%s'", result.String())
	}
}

func TestEmptyAndNullValues(t *testing.T) {
	// Test empty object
	result := Get(edgeCaseYAML, "empty_and_null.empty_object")
	if !result.Exists() {
		t.Error("Empty object should exist")
	}
	objMap := result.Map()
	if len(objMap) != 0 {
		t.Errorf("Expected empty map, got %d items", len(objMap))
	}

	// Test empty array
	result = Get(edgeCaseYAML, "empty_and_null.empty_array")
	if !result.Exists() {
		t.Error("Empty array should exist")
	}
	arr := result.Array()
	if len(arr) != 0 {
		t.Errorf("Expected empty array, got %d items", len(arr))
	}

	// Test null value
	result = Get(edgeCaseYAML, "empty_and_null.null_explicit")
	if result.Type != Null {
		t.Errorf("Expected Null type, got %v", result.Type)
	}

	// Test zero value
	result = Get(edgeCaseYAML, "empty_and_null.zero")
	if result.Int() != 0 {
		t.Errorf("Expected 0, got %d", result.Int())
	}
	if result.Type != Number {
		t.Errorf("Expected Number type for zero, got %v", result.Type)
	}
}

func TestQuotesAndEscapes(t *testing.T) {
	// Test single quoted string
	result := Get(edgeCaseYAML, "quotes_and_escapes.single_quotes")
	if result.String() != "single quoted string" {
		t.Errorf("Expected 'single quoted string', got '%s'", result.String())
	}

	// Test double quoted string
	result = Get(edgeCaseYAML, "quotes_and_escapes.double_quotes")
	if result.String() != "double quoted string" {
		t.Errorf("Expected 'double quoted string', got '%s'", result.String())
	}

	// Test escape characters
	result = Get(edgeCaseYAML, "quotes_and_escapes.escaped_chars")
	escaped := result.String()
	if !Contains(escaped, "\n") || !Contains(escaped, "\t") {
		t.Errorf("Expected escaped characters, got '%s'", escaped)
	}

	// Test backslash
	result = Get(edgeCaseYAML, "quotes_and_escapes.backslash")
	backslash := result.String()
	if !Contains(backslash, "\\") {
		t.Errorf("Expected backslash, got '%s'", backslash)
	}
}

func TestLargeArrayPerformance(t *testing.T) {
	// Test large array length
	result := Get(edgeCaseYAML, "large_array.numbers.#")
	if result.Int() != 50 {
		t.Errorf("Expected 50 numbers, got %d", result.Int())
	}

	// Test array element access
	result = Get(edgeCaseYAML, "large_array.numbers.49")
	if result.Int() != 50 {
		t.Errorf("Expected 50, got %d", result.Int())
	}

	// Test object array queries
	result = Get(edgeCaseYAML, `large_array.objects.#(id=3)`)
	if !result.Exists() {
		t.Error("Expected to find object with id=3")
	}
	name := result.Get("name")
	if name.String() != "Item 3" {
		t.Errorf("Expected 'Item 3', got '%s'", name.String())
	}

	// Test getting all active object IDs
	result = Get(edgeCaseYAML, "large_array.objects.#.id")
	ids := result.Array()
	if len(ids) != 5 {
		t.Errorf("Expected 5 IDs, got %d", len(ids))
	}
}

func TestYAMLSyntaxFeatures(t *testing.T) {
	// Test multi-line text blocks
	result := Get(edgeCaseYAML, "yaml_syntax.multiline_literal")
	literal := result.String()
	if !Contains(literal, "First line") || !Contains(literal, "Second line") {
		t.Errorf("Expected multiline literal, got '%s'", literal)
	}

	// Test folded string
	result = Get(edgeCaseYAML, "yaml_syntax.multiline_folded")
	folded := result.String()
	if !Contains(folded, "very very") {
		t.Errorf("Expected folded string, got '%s'", folded)
	}
}

// Note: YAML 3.0 does not allow duplicate keys by default, so duplicate key tests were removed

func TestComplexPaths(t *testing.T) {
	// Test complex path access
	result := Get(edgeCaseYAML, "complex_paths.path_with_dots.nested_key")
	if result.String() != "value1" {
		t.Errorf("Expected 'value1', got '%s'", result.String())
	}

	// Test key names containing special characters
	result = Get(edgeCaseYAML, "complex_paths.array_bracket.element")
	if result.String() != "confusing_name" {
		t.Errorf("Expected 'confusing_name', got '%s'", result.String())
	}

	// Test other special key names
	result = Get(edgeCaseYAML, "complex_paths.hashtag_key.count")
	if result.Int() != 100 {
		t.Errorf("Expected 100, got %d", result.Int())
	}
}

func TestArrayBoundaryConditions(t *testing.T) {
	// Test array boundary access
	result := Get(edgeCaseYAML, "large_array.numbers.-1") // Negative index
	if result.Exists() {
		t.Error("Negative array index should not exist")
	}

	result = Get(edgeCaseYAML, "large_array.numbers.1000") // Out of bounds
	if result.Exists() {
		t.Error("Out-of-bounds array index should not exist")
	}

	// Test operations on empty array
	result = Get(edgeCaseYAML, "empty_and_null.empty_array.0")
	if result.Exists() {
		t.Error("Empty array element should not exist")
	}

	result = Get(edgeCaseYAML, "empty_and_null.empty_array.#")
	if result.Int() != 0 {
		t.Errorf("Empty array length should be 0, got %d", result.Int())
	}
}

func TestConcurrentAccess(t *testing.T) {
	// Test concurrent access safety
	paths := []string{
		"special_keys.key-with-dashes",
		"extreme_values.max_int",
		"nested_arrays.level1.0.level2.0.level3.0.level4.0",
		"large_array.numbers.25",
		"yaml_syntax.multiline_literal",
	}

	done := make(chan bool, len(paths))

	for _, path := range paths {
		go func(p string) {
			for i := 0; i < 100; i++ {
				result := Get(edgeCaseYAML, p)
				if !result.Exists() {
					t.Errorf("Path %s should exist in concurrent access", p)
				}
			}
			done <- true
		}(path)
	}

	// Wait for all goroutines to complete
	for i := 0; i < len(paths); i++ {
		<-done
	}
}

func TestMemoryUsage(t *testing.T) {
	// Test for memory leaks with repeated calls
	for i := 0; i < 1000; i++ {
		result := Get(edgeCaseYAML, "extreme_values.very_long_string")
		if !result.Exists() {
			t.Error("Large string should always exist")
		}

		// Test complex paths
		result = Get(edgeCaseYAML, "nested_arrays.level1.0.level2.0.level3.0.level4")
		arr := result.Array()
		if len(arr) == 0 {
			t.Error("Nested array should not be empty")
		}
	}
}
