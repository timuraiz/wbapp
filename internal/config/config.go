package config

import (
	"gopkg.in/yaml.v2"
	"os"
)

type Config struct {
	WeeklyReportConfig `yaml:"weekly-report-config"`
}

type WeeklyReportConfig struct {
	Host string `yaml:"host"`
}

type Method struct {
}

func ParseConfig(path string) (*Config, error) {
	fileData, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var data Config
	err = yaml.Unmarshal(fileData, &data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}
