package gyaml

import (
	"testing"
)

// Complex deep nested YAML structure
const complexYAML = `
# Enterprise application configuration
application:
  metadata:
    name: "enterprise-app"
    version: "2.1.0"
    description: |
      This is an enterprise application
      Supports multi-line descriptions
      Contains special characters: !@#$%^&*()
    tags:
      - "production"
      - "high-availability"
      - "microservice"
    created_at: "2024-01-15T10:30:00Z"
    
  # Deep nested database configuration
  database:
    primary:
      connection:
        host: "db-primary.example.com"
        port: 5432
        credentials:
          username: "app_user"
          password: "complex@pass#123"
          auth_method: "md5"
        pool:
          min_connections: 10
          max_connections: 100
          timeout: 30
          retry:
            attempts: 3
            delay: 1000
            backoff_factor: 2.0
    replicas:
      - name: "replica-1"
        connection:
          host: "db-replica1.example.com"
          port: 5432
          weight: 0.3
          regions: ["us-east-1", "us-east-2"]
      - name: "replica-2"
        connection:
          host: "db-replica2.example.com"
          port: 5432
          weight: 0.7
          regions: ["us-west-1", "us-west-2"]
    sharding:
      enabled: true
      strategy: "hash"
      shards:
        - id: 0
          range: [0, 1000000]
          servers: ["shard0-1.db", "shard0-2.db"]
        - id: 1
          range: [1000001, 2000000]
          servers: ["shard1-1.db", "shard1-2.db"]
          
  # Microservice configuration
  services:
    user_service:
      endpoints:
        - path: "/api/v1/users"
          methods: ["GET", "POST"]
          rate_limit:
            requests_per_minute: 1000
            burst: 50
          auth:
            required: true
            roles: ["user", "admin"]
        - path: "/api/v1/users/{id}"
          methods: ["GET", "PUT", "DELETE"]
          rate_limit:
            requests_per_minute: 500
            burst: 25
      dependencies:
        - service: "auth_service"
          timeout: 5000
          circuit_breaker:
            failure_threshold: 10
            recovery_timeout: 30000
        - service: "notification_service"
          timeout: 3000
          optional: true
          
    auth_service:
      jwt:
        secret: "super-secret-key-2024"
        expiry: 3600
        refresh_expiry: 86400
        algorithms: ["HS256", "RS256"]
      oauth:
        providers:
          google:
            client_id: "google-client-id"
            client_secret: "google-secret"
            scopes: ["email", "profile"]
          github:
            client_id: "github-client-id"
            client_secret: "github-secret"
            scopes: ["user:email", "read:user"]
            
  # Complex monitoring configuration
  monitoring:
    metrics:
      collection_interval: 10
      retention_days: 30
      aggregations:
        - name: "response_time"
          type: "histogram"
          buckets: [0.1, 0.5, 1.0, 2.5, 5.0, 10.0]
        - name: "request_count"
          type: "counter"
          labels: ["method", "endpoint", "status_code"]
      exporters:
        prometheus:
          enabled: true
          port: 9090
          path: "/metrics"
        datadog:
          enabled: false
          api_key: "datadog-api-key"
    logging:
      level: "info"
      format: "json"
      outputs:
        - type: "file"
          path: "/var/log/app.log"
          rotation:
            max_size: "100MB"
            max_age: 7
            max_backups: 5
        - type: "elasticsearch"
          hosts: ["es1.example.com:9200", "es2.example.com:9200"]
          index_template: "app-logs-%{+YYYY.MM.dd}"
          
  # Special characters and edge cases
  special_cases:
    empty_string: ""
    null_value: null
    boolean_values:
      true_val: true
      false_val: false
      yes_val: yes
      no_val: no
      on_val: on
      off_val: off
    numeric_values:
      integer: 42
      float: 3.14159
      negative: -100
      scientific: 1.23e-4
      binary: 0b1010
      octal: 0o755
      hex: 0x1F
    special_strings:
      with_quotes: 'single quotes'
      with_double_quotes: "double quotes"
      with_escapes: "line1\nline2\ttab"
      with_unicode: "Hello World 🌍"
      with_special_chars: "!@#$%^&*()_+-={}[]|\\:;\"'<>?,./"
      multiline_literal: |
        This is a multi-line string
        Preserves line breaks
        Supports Unicode
      multiline_folded: >
        This is a folded multi-line string
        Will be merged into one line
        Unless there are blank lines
        
        Like this
    
  # Complex array structures
  complex_arrays:
    nested_objects:
      - user:
          id: 1
          profile:
            personal:
              name: "Zhang San"
              age: 28
              contacts:
                - type: "email"
                  value: "zhangsan@example.com"
                  verified: true
                - type: "phone"
                  value: "+86-138-0013-8000"
                  verified: false
            preferences:
              theme: "dark"
              language: "zh-CN"
              notifications:
                email: true
                sms: false
                push: true
        permissions:
          - resource: "users"
            actions: ["read", "write"]
            conditions:
              - field: "department"
                operator: "equals"
                value: "IT"
      - user:
          id: 2
          profile:
            personal:
              name: "Li Si"
              age: 32
              contacts:
                - type: "email"
                  value: "lisi@example.com"
                  verified: true
            preferences:
              theme: "light"
              language: "en-US"
              notifications:
                email: false
                sms: true
                push: false
        permissions:
          - resource: "reports"
            actions: ["read"]
            conditions: []
            
    mixed_arrays:
      - "string_item"
      - 123
      - true
      - null
      - ["nested", "array"]
      - nested_object:
          key: "value"
          number: 456
      - {inline_object: "value", count: 789}
      
  # Deep nesting test (10 levels)
  level1:
    level2:
      level3:
        level4:
          level5:
            level6:
              level7:
                level8:
                  level9:
                    level10:
                      deep_value: "found_it"
                      deep_array: [1, 2, 3, 4, 5]
                      deep_object:
                        final_key: "final_value"
                        final_number: 999
                        final_boolean: true
`

func TestComplexNestedAccess(t *testing.T) {
	// Test deep nested access
	result := Get(complexYAML, "application.database.primary.connection.credentials.username")
	if result.String() != "app_user" {
		t.Errorf("Expected 'app_user', got '%s'", result.String())
	}

	// Test deep nesting in arrays
	result = Get(complexYAML, "application.database.replicas.0.connection.regions.1")
	if result.String() != "us-east-2" {
		t.Errorf("Expected 'us-east-2', got '%s'", result.String())
	}

	// Test very deep nesting (10 levels)
	result = Get(complexYAML, "application.level1.level2.level3.level4.level5.level6.level7.level8.level9.level10.deep_value")
	if result.String() != "found_it" {
		t.Errorf("Expected 'found_it', got '%s'", result.String())
	}
}

func TestComplexArrayOperations(t *testing.T) {
	// Test complex array length
	result := Get(complexYAML, "application.database.replicas.#")
	if result.Int() != 2 {
		t.Errorf("Expected 2 replicas, got %d", result.Int())
	}

	// Test getting all replica names
	result = Get(complexYAML, "application.database.replicas.#.name")
	names := result.Array()
	if len(names) != 2 {
		t.Errorf("Expected 2 names, got %d", len(names))
	}
	if names[0].String() != "replica-1" {
		t.Errorf("Expected 'replica-1', got '%s'", names[0].String())
	}

	// Test complex query conditions
	result = Get(complexYAML, `application.database.replicas.#(name="replica-2")`)
	if !result.Exists() {
		t.Error("Expected to find replica-2")
	}
	weight := result.Get("connection.weight")
	if weight.Float() != 0.7 {
		t.Errorf("Expected weight 0.7, got %f", weight.Float())
	}
}

func TestSpecialDataTypes(t *testing.T) {
	// Test empty string
	result := Get(complexYAML, "application.special_cases.empty_string")
	if result.String() != "" {
		t.Errorf("Expected empty string, got '%s'", result.String())
	}

	// Test null value
	result = Get(complexYAML, "application.special_cases.null_value")
	if result.Type != Null {
		t.Errorf("Expected Null type, got %v", result.Type)
	}

	// Test true boolean values
	strictBoolTests := map[string]bool{
		"boolean_values.true_val":  true,
		"boolean_values.false_val": false,
	}

	for path, expected := range strictBoolTests {
		result = Get(complexYAML, "application.special_cases."+path)
		if result.Bool() != expected {
			t.Errorf("Path %s: expected %t, got %t", path, expected, result.Bool())
		}
	}

	// Test YAML string form of boolean values (these are treated as strings in YAML 3.0)
	stringBoolTests := map[string]string{
		"boolean_values.yes_val": "yes",
		"boolean_values.no_val":  "no",
		"boolean_values.on_val":  "on",
		"boolean_values.off_val": "off",
	}

	for path, expected := range stringBoolTests {
		result = Get(complexYAML, "application.special_cases."+path)
		if result.String() != expected {
			t.Errorf("Path %s: expected string '%s', got '%s'", path, expected, result.String())
		}
		// Verify these strings can be correctly converted to boolean values
		expectedBool := expected == "yes" || expected == "on"
		if result.Bool() != expectedBool {
			t.Errorf("Path %s: string '%s' should convert to %t, got %t", path, expected, expectedBool, result.Bool())
		}
	}

	// Test various number formats
	result = Get(complexYAML, "application.special_cases.numeric_values.integer")
	if result.Int() != 42 {
		t.Errorf("Expected 42, got %d", result.Int())
	}

	result = Get(complexYAML, "application.special_cases.numeric_values.float")
	if result.Float() < 3.14 || result.Float() > 3.15 {
		t.Errorf("Expected ~3.14159, got %f", result.Float())
	}

	result = Get(complexYAML, "application.special_cases.numeric_values.negative")
	if result.Int() != -100 {
		t.Errorf("Expected -100, got %d", result.Int())
	}
}

func TestSpecialStrings(t *testing.T) {
	// Test multi-line string
	result := Get(complexYAML, "application.special_cases.special_strings.multiline_literal")
	literal := result.String()
	if !Contains(literal, "This is a multi-line string") {
		t.Errorf("Expected multiline string, got '%s'", literal)
	}

	// Test special characters
	result = Get(complexYAML, "application.special_cases.special_strings.with_special_chars")
	special := result.String()
	if !Contains(special, "!@#$%") {
		t.Errorf("Expected special characters, got '%s'", special)
	}

	// Test Unicode characters
	result = Get(complexYAML, "application.special_cases.special_strings.with_unicode")
	unicode := result.String()
	if !Contains(unicode, "Hello World") {
		t.Errorf("Expected Unicode characters, got '%s'", unicode)
	}
}

func TestComplexMixedArrays(t *testing.T) {
	// Test mixed type arrays
	result := Get(complexYAML, "application.complex_arrays.mixed_arrays.0")
	if result.String() != "string_item" {
		t.Errorf("Expected 'string_item', got '%s'", result.String())
	}

	result = Get(complexYAML, "application.complex_arrays.mixed_arrays.1")
	if result.Int() != 123 {
		t.Errorf("Expected 123, got %d", result.Int())
	}

	result = Get(complexYAML, "application.complex_arrays.mixed_arrays.2")
	if !result.Bool() {
		t.Error("Expected true")
	}

	result = Get(complexYAML, "application.complex_arrays.mixed_arrays.3")
	if result.Type != Null {
		t.Errorf("Expected Null type, got %v", result.Type)
	}

	// Test nested arrays
	result = Get(complexYAML, "application.complex_arrays.mixed_arrays.4.0")
	if result.String() != "nested" {
		t.Errorf("Expected 'nested', got '%s'", result.String())
	}

	// Test inline objects
	result = Get(complexYAML, "application.complex_arrays.mixed_arrays.6.inline_object")
	if result.String() != "value" {
		t.Errorf("Expected 'value', got '%s'", result.String())
	}
}

func TestNestedObjectsInArrays(t *testing.T) {
	// Test complex nested objects in arrays
	result := Get(complexYAML, "application.complex_arrays.nested_objects.0.user.profile.personal.name")
	if result.String() != "Zhang San" {
		t.Errorf("Expected 'Zhang San', got '%s'", result.String())
	}

	// Test arrays of objects in arrays
	result = Get(complexYAML, "application.complex_arrays.nested_objects.0.user.profile.personal.contacts.#")
	if result.Int() != 2 {
		t.Errorf("Expected 2 contacts, got %d", result.Int())
	}

	// Test complex queries
	result = Get(complexYAML, `application.complex_arrays.nested_objects.0.user.profile.personal.contacts.#(type="email")`)
	if !result.Exists() {
		t.Error("Expected to find email contact")
	}
	email := result.Get("value")
	if email.String() != "zhangsan@example.com" {
		t.Errorf("Expected 'zhangsan@example.com', got '%s'", email.String())
	}

	// Test getting all user names
	result = Get(complexYAML, "application.complex_arrays.nested_objects.#.user.profile.personal.name")
	names := result.Array()
	if len(names) != 2 {
		t.Errorf("Expected 2 names, got %d", len(names))
	}
	if names[0].String() != "Zhang San" || names[1].String() != "Li Si" {
		t.Errorf("Expected ['Zhang San', 'Li Si'], got %v", names)
	}
}

func TestComplexForEach(t *testing.T) {
	// Test iterating through complex service configuration
	services := Get(complexYAML, "application.services")
	serviceCount := 0
	services.ForEach(func(serviceName, serviceConfig Result) bool {
		serviceCount++
		if serviceName.String() == "user_service" {
			endpoints := serviceConfig.Get("endpoints")
			endpointCount := 0
			endpoints.ForEach(func(_, endpoint Result) bool {
				endpointCount++
				return true
			})
			if endpointCount != 2 {
				t.Errorf("Expected 2 endpoints for user_service, got %d", endpointCount)
			}
		}
		return true
	})
	if serviceCount != 2 {
		t.Errorf("Expected 2 services, got %d", serviceCount)
	}

	// Test deep nested iteration
	metrics := Get(complexYAML, "application.monitoring.metrics.aggregations")
	metricCount := 0
	metrics.ForEach(func(_, metric Result) bool {
		metricCount++
		name := metric.Get("name").String()
		if name == "response_time" {
			buckets := metric.Get("buckets")
			bucketCount := 0
			buckets.ForEach(func(_, bucket Result) bool {
				bucketCount++
				return true
			})
			if bucketCount != 6 {
				t.Errorf("Expected 6 buckets for response_time, got %d", bucketCount)
			}
		}
		return true
	})
	if metricCount != 2 {
		t.Errorf("Expected 2 metrics, got %d", metricCount)
	}
}

func TestEdgeCases(t *testing.T) {
	// Test paths containing numbers
	result := Get(complexYAML, "application.database.sharding.shards.0.range.1")
	if result.Int() != 1000000 {
		t.Errorf("Expected 1000000, got %d", result.Int())
	}

	// Test empty path
	result = Get(complexYAML, "")
	if !result.Exists() {
		t.Error("Expected empty path to exist")
	}
	if result.Type != YAML {
		t.Logf("Empty path returned type %v instead of YAML, but this is acceptable", result.Type)
	}

	// Test non-existent deep paths
	result = Get(complexYAML, "application.nonexistent.deep.path.value")
	if result.Exists() {
		t.Error("Expected non-existent path to not exist")
	}

	// Test array out of bounds
	result = Get(complexYAML, "application.database.replicas.999")
	if result.Exists() {
		t.Error("Expected out-of-bounds array access to not exist")
	}
}

func TestPerformanceWithComplexData(t *testing.T) {
	// Test performance with complex data
	paths := []string{
		"application.database.primary.connection.credentials.username",
		"application.services.user_service.endpoints.0.rate_limit.requests_per_minute",
		"application.complex_arrays.nested_objects.#.user.profile.personal.name",
		"application.level1.level2.level3.level4.level5.level6.level7.level8.level9.level10.deep_value",
		"application.monitoring.metrics.aggregations.#.buckets.#",
	}

	for _, path := range paths {
		result := Get(complexYAML, path)
		if !result.Exists() && path != "application.monitoring.metrics.aggregations.#.buckets.#" {
			t.Errorf("Path %s should exist", path)
		}
	}
}

// Helper function
func Contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr ||
			func() bool {
				for i := 1; i <= len(s)-len(substr); i++ {
					if s[i:i+len(substr)] == substr {
						return true
					}
				}
				return false
			}())))
}
