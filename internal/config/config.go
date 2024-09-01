package config

import (
	"gopkg.in/yaml.v2"
	"io"
	"time"
)

type Config struct {
	WeeklyReportConfig `yaml:"weekly-report-config"`
}

type WeeklyReportConfig struct {
	Host    string   `yaml:"host"`
	Methods []Method `yaml:"methods"`
}

type Method struct {
	Method  string        `yaml:"method"`
	Route   string        `yaml:"route"`
	Timeout time.Duration `yaml:"timeout"`
}

func ParseConfig(buffer io.Reader) (*Config, error) {
	bytes, err := io.ReadAll(buffer)
	if err != nil {
		return nil, err
	}

	var data Config
	err = yaml.Unmarshal(bytes, &data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}
