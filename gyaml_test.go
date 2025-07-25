package gyaml

import (
	"testing"
)

const testYAML = `
name:
  first: "Tom"
  last: "Anderson"
age: 37
children:
  - "Sara"
  - "Alex"
  - "Jack"
fav_movie: "Deer Hunter"
friends:
  - first: "Dale"
    last: "Murphy"
    age: 44
    hobbies: ["golf", "tennis"]
  - first: "Roger"
    last: "Craig"
    age: 68
    hobbies: ["fishing", "cooking"]
  - first: "Jane"
    last: "Murphy"
    age: 47
    hobbies: ["reading", "gardening"]
`

func TestGet(t *testing.T) {
	// Test basic path access
	result := Get(testYAML, "name.first")
	if !result.Exists() {
		t.Error("Expected result to exist")
	}
	if result.String() != "Tom" {
		t.Errorf("Expected 'Tom', got '%s'", result.String())
	}

	// Test number access
	result = Get(testYAML, "age")
	if result.Int() != 37 {
		t.Errorf("Expected 37, got %d", result.Int())
	}

	// Test array access
	result = Get(testYAML, "children.0")
	if result.String() != "Sara" {
		t.Errorf("Expected 'Sara', got '%s'", result.String())
	}

	// Test array length
	result = Get(testYAML, "children.#")
	if result.Int() != 3 {
		t.Errorf("Expected 3, got %d", result.Int())
	}

	// Test non-existent path
	result = Get(testYAML, "nonexistent")
	if result.Exists() {
		t.Error("Expected result to not exist")
	}
}

func TestArrayOperations(t *testing.T) {
	// Test getting all first names from friends
	result := Get(testYAML, "friends.#.first")
	arr := result.Array()
	if len(arr) != 3 {
		t.Errorf("Expected 3 items, got %d", len(arr))
	}
	if arr[0].String() != "Dale" {
		t.Errorf("Expected 'Dale', got '%s'", arr[0].String())
	}
}

func TestArrayQuery(t *testing.T) {
	// Test finding friend by name
	result := Get(testYAML, `friends.#(first="Roger")`)
	if !result.Exists() {
		t.Error("Expected result to exist")
	}

	friendResult := result.Get("last")
	if friendResult.String() != "Craig" {
		t.Errorf("Expected 'Craig', got '%s'", friendResult.String())
	}
}

func TestParse(t *testing.T) {
	result := Parse(testYAML)
	if !result.Exists() {
		t.Error("Expected result to exist")
	}

	// Test accessing through parsed result
	name := result.Get("name.first")
	if name.String() != "Tom" {
		t.Errorf("Expected 'Tom', got '%s'", name.String())
	}
}

func TestValid(t *testing.T) {
	if !Valid(testYAML) {
		t.Error("Expected YAML to be valid")
	}

	if Valid("invalid: yaml: [") {
		t.Error("Expected invalid YAML to be invalid")
	}
}

func TestForEach(t *testing.T) {
	result := Get(testYAML, "name")
	count := 0
	result.ForEach(func(key, value Result) bool {
		count++
		return true
	})
	if count != 2 {
		t.Errorf("Expected 2 items, got %d", count)
	}

	// Test array iteration
	result = Get(testYAML, "children")
	count = 0
	result.ForEach(func(key, value Result) bool {
		count++
		return true
	})
	if count != 3 {
		t.Errorf("Expected 3 items, got %d", count)
	}
}

func TestTypes(t *testing.T) {
	// Test string
	result := Get(testYAML, "name.first")
	if result.Type != String {
		t.Error("Expected String type")
	}

	// Test number
	result = Get(testYAML, "age")
	if result.Type != Number {
		t.Error("Expected Number type")
	}

	// Test array
	result = Get(testYAML, "children")
	if result.Type != YAML {
		t.Error("Expected YAML type for array")
	}
}

func TestGetBytes(t *testing.T) {
	yamlBytes := []byte(testYAML)
	result := GetBytes(yamlBytes, "name.first")
	if result.String() != "Tom" {
		t.Errorf("Expected 'Tom', got '%s'", result.String())
	}
}
