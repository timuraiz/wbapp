package config

import (
	"strings"
	"testing"
	"time"
)

func TestParseConfig(t *testing.T) {
	yamlInput := `
weekly-report-config:
  host: "url"
  methods:
    - method: "categories"
      route: "/api/v1/documents/categories"
      timeout: 10s
    - method: "list"
      route: "/api/v1/documents/list"
      timeout: 10s
    - method: "download"
      route: "/api/v1/documents/download"
      timeout: 10s
`

	// Create a new reader from the YAML input string
	r := strings.NewReader(yamlInput)

	// Parse the config
	config, err := ParseConfig(r)
	if err != nil {
		t.Fatalf("Failed to parse config: %v", err)
	}

	// Define the expected timeout for the methods
	expectedTimeout := 10 * time.Second

	// Check overall WeeklyReportConfig.Host
	if config.WeeklyReportConfig.Host != "url" {
		t.Errorf("Expected host to be 'url', but got %s", config.WeeklyReportConfig.Host)
	}

	// Check each method's properties
	expectedMethods := []struct {
		Method  string
		Route   string
		Timeout time.Duration
	}{
		{"categories", "/api/v1/documents/categories", expectedTimeout},
		{"list", "/api/v1/documents/list", expectedTimeout},
		{"download", "/api/v1/documents/download", expectedTimeout},
	}

	for i, method := range config.WeeklyReportConfig.Methods {
		if method.Method != expectedMethods[i].Method {
			t.Errorf("Expected method %d to be '%s', but got '%s'", i, expectedMethods[i].Method, method.Method)
		}
		if method.Route != expectedMethods[i].Route {
			t.Errorf("Expected route %d to be '%s', but got '%s'", i, expectedMethods[i].Route, method.Route)
		}
		if method.Timeout != expectedMethods[i].Timeout {
			t.Errorf("Expected timeout %d to be %v, but got %v", i, expectedMethods[i].Timeout, method.Timeout)
		}
	}
}
