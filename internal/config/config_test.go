package config

import (
	"testing"
	"time"
)

func TestParseConfig(t *testing.T) {
	config := ParseConfig("../../configs/config.yaml")
	expected := 10 * time.Second
	if config.WeeklyReportConfig.Timeout != expected {
		t.Fatalf("Expected %d but got %d", expected, config.WeeklyReportConfig.Timeout)
	}
}
