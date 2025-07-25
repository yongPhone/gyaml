package gyaml

import (
	"testing"
)

const benchmarkYAML = `
users:
  - id: 1
    name: "Alice Johnson"
    email: "alice@example.com"
    profile:
      age: 28
      city: "New York"
      hobbies: ["reading", "swimming", "coding"]
      settings:
        theme: "dark"
        notifications: true
  - id: 2
    name: "Bob Smith"
    email: "bob@example.com"
    profile:
      age: 34
      city: "San Francisco"
      hobbies: ["hiking", "photography"]
      settings:
        theme: "light"
        notifications: false
  - id: 3
    name: "Charlie Brown"
    email: "charlie@example.com"
    profile:
      age: 22
      city: "Seattle"
      hobbies: ["gaming", "music", "cooking"]
      settings:
        theme: "auto"
        notifications: true
config:
  database:
    host: "localhost"
    port: 5432
    ssl: true
  server:
    host: "0.0.0.0"
    port: 8080
    workers: 4
`

func BenchmarkGet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Get(benchmarkYAML, "users.0.name")
	}
}

func BenchmarkGetNested(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Get(benchmarkYAML, "users.0.profile.settings.theme")
	}
}

func BenchmarkGetArray(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Get(benchmarkYAML, "users.0.profile.hobbies.1")
	}
}

func BenchmarkGetArrayLength(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Get(benchmarkYAML, "users.#")
	}
}

func BenchmarkGetArrayOperation(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Get(benchmarkYAML, "users.#.name")
	}
}

func BenchmarkGetArrayQuery(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Get(benchmarkYAML, `users.#(id=2)`)
	}
}

func BenchmarkParse(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Parse(benchmarkYAML)
	}
}

func BenchmarkValid(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Valid(benchmarkYAML)
	}
}

func BenchmarkForEach(b *testing.B) {
	result := Get(benchmarkYAML, "users")
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		result.ForEach(func(key, value Result) bool {
			return true
		})
	}
}

func BenchmarkGetMultiple(b *testing.B) {
	paths := []string{
		"users.0.name",
		"users.1.email",
		"users.2.profile.age",
		"config.database.host",
		"config.server.port",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, path := range paths {
			Get(benchmarkYAML, path)
		}
	}
}
