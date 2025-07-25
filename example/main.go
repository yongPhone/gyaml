package main

import (
	"fmt"

	"github.com/yongPhone/gyaml"
)

func main() {
	// Example YAML data
	yamlData := `
# Application configuration
app:
  name: "My Application"
  version: "1.0.0"
  debug: true
  
# Database configuration
database:
  host: "localhost"
  port: 5432
  username: "admin"
  password: "secret"
  
# Server list
servers:
  - name: "web1"
    ip: "192.168.1.10"
    roles: ["web", "api"]
  - name: "web2"  
    ip: "192.168.1.11"
    roles: ["web"]
  - name: "db1"
    ip: "192.168.1.20"
    roles: ["database"]

# Feature flags
features:
  new_ui: true
  analytics: false
  beta_feature: true
`

	fmt.Println("=== GYAML Usage Examples ===")

	// 1. Basic field access
	fmt.Println("1. Basic field access:")
	appName := gyaml.Get(yamlData, "app.name")
	fmt.Printf("Application name: %s\n", appName.String())

	appVersion := gyaml.Get(yamlData, "app.version")
	fmt.Printf("Version: %s\n", appVersion.String())

	isDebug := gyaml.Get(yamlData, "app.debug")
	fmt.Printf("Debug mode: %t\n\n", isDebug.Bool())

	// 2. Array operations
	fmt.Println("2. Array operations:")
	serverCount := gyaml.Get(yamlData, "servers.#")
	fmt.Printf("Total servers: %d\n", serverCount.Int())

	firstServer := gyaml.Get(yamlData, "servers.0.name")
	fmt.Printf("First server: %s\n", firstServer.String())

	// Get all server names
	fmt.Println("All server names:")
	serverNames := gyaml.Get(yamlData, "servers.#.name")
	for _, name := range serverNames.Array() {
		fmt.Printf("  - %s\n", name.String())
	}
	fmt.Println()

	// 3. Array queries
	fmt.Println("3. Array queries:")
	webServer := gyaml.Get(yamlData, `servers.#(name="web1")`)
	if webServer.Exists() {
		ip := webServer.Get("ip").String()
		fmt.Printf("web1 server IP: %s\n", ip)
	}
	fmt.Println()

	// 4. Iterate through objects
	fmt.Println("4. Iterate through database config:")
	dbConfig := gyaml.Get(yamlData, "database")
	dbConfig.ForEach(func(key, value gyaml.Result) bool {
		fmt.Printf("  %s: %s\n", key.String(), value.String())
		return true
	})
	fmt.Println()

	// 5. Iterate through arrays
	fmt.Println("5. Iterate through server list:")
	servers := gyaml.Get(yamlData, "servers")
	servers.ForEach(func(index, server gyaml.Result) bool {
		name := server.Get("name").String()
		ip := server.Get("ip").String()
		fmt.Printf("  Server %d: %s (%s)\n", int(index.Num), name, ip)
		return true
	})
	fmt.Println()

	// 6. Check field existence
	fmt.Println("6. Check field existence:")
	if gyaml.Get(yamlData, "app.timeout").Exists() {
		fmt.Println("Application timeout config exists")
	} else {
		fmt.Println("Application timeout config does not exist")
	}

	if gyaml.Get(yamlData, "features.new_ui").Exists() {
		enabled := gyaml.Get(yamlData, "features.new_ui").Bool()
		fmt.Printf("New UI feature: %s\n", map[bool]string{true: "enabled", false: "disabled"}[enabled])
	}
	fmt.Println()

	// 7. Data type conversion
	fmt.Println("7. Data type conversion:")
	dbPort := gyaml.Get(yamlData, "database.port")
	fmt.Printf("Database port (integer): %d\n", dbPort.Int())
	fmt.Printf("Database port (string): %s\n", dbPort.String())
	fmt.Printf("Database port (float): %.1f\n", dbPort.Float())
	fmt.Println()

	// 8. YAML validation
	fmt.Println("8. YAML validation:")
	if gyaml.Valid(yamlData) {
		fmt.Println("YAML format is valid")
	} else {
		fmt.Println("YAML format is invalid")
	}

	invalidYAML := `
invalid: yaml
  missing: quotes
    bad: indent
`
	fmt.Printf("Invalid YAML check: %t\n", gyaml.Valid(invalidYAML))
}
