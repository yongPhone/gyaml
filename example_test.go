package gyaml_test

import (
	"fmt"

	"github.com/yongPhone/gyaml"
)

func ExampleGet() {
	yaml := `
name:
  first: "Janet"
  last: "Prichard"
age: 47
`
	value := gyaml.Get(yaml, "name.last")
	fmt.Println(value.String())
	// Output: Prichard
}

func ExampleGet_array() {
	yaml := `
children:
  - "Sara"
  - "Alex"
  - "Jack"
`
	// Get array length
	length := gyaml.Get(yaml, "children.#")
	fmt.Println("Length:", length.Int())

	// Get specific element
	first := gyaml.Get(yaml, "children.0")
	fmt.Println("First child:", first.String())

	// Output:
	// Length: 3
	// First child: Sara
}

func ExampleGet_query() {
	yaml := `
programmers:
  - firstName: "Janet"
    lastName: "McLaughlin"
    language: "Go"
  - firstName: "Elliotte"
    lastName: "Hunter"
    language: "Java"
  - firstName: "Jason"
    lastName: "Harold"
    language: "Python"
`

	// Find programmer by language
	result := gyaml.Get(yaml, `programmers.#(language="Java")`)
	fmt.Println("Java programmer:", result.Get("firstName").String())

	// Get all first names
	names := gyaml.Get(yaml, "programmers.#.firstName")
	for _, name := range names.Array() {
		fmt.Println("Name:", name.String())
	}

	// Output:
	// Java programmer: Elliotte
	// Name: Janet
	// Name: Elliotte
	// Name: Jason
}

func ExampleResult_ForEach() {
	yaml := `
users:
  - name: "Alice"
    role: "admin"
  - name: "Bob"
    role: "user"
`

	result := gyaml.Get(yaml, "users")
	count := 0
	result.ForEach(func(key, value gyaml.Result) bool {
		count++
		name := value.Get("name").String()
		role := value.Get("role").String()
		fmt.Printf("User %s: %s (%s)\n", key.String(), name, role)
		return true // continue iteration
	})
	fmt.Printf("Total users: %d\n", count)

	// Output:
	// User 0: Alice (admin)
	// User 1: Bob (user)
	// Total users: 2
}

func ExampleParse() {
	yaml := `
config:
  debug: true
  timeout: 30
  features:
    - "feature1"
    - "feature2"
`

	result := gyaml.Parse(yaml)

	// Check if debug is enabled
	if result.Get("config.debug").Bool() {
		fmt.Println("Debug mode enabled")
	}

	// Get timeout value
	timeout := result.Get("config.timeout").Int()
	fmt.Printf("Timeout: %d seconds\n", timeout)

	// Check number of features
	featureCount := result.Get("config.features.#").Int()
	fmt.Printf("Features count: %d\n", featureCount)

	// Output:
	// Debug mode enabled
	// Timeout: 30 seconds
	// Features count: 2
}

func ExampleValid() {
	validYAML := `
name: "John"
age: 30
`

	invalidYAML := `
name: "John"
age: 30
  invalid: syntax
`

	fmt.Println("Valid YAML:", gyaml.Valid(validYAML))
	fmt.Println("Invalid YAML:", gyaml.Valid(invalidYAML))

	// Output:
	// Valid YAML: true
	// Invalid YAML: false
}
