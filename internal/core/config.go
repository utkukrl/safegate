package core

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Proxy struct {
		Target string `yaml:"target"`
		Port   string `yaml:"port"`
	} `yaml:"proxy"`

	Firewall struct {
		IPWhitelist []string `yaml:"ip_whitelist"`
	} `yaml:"firewall"`

	RateLimit struct {
		Enabled bool   `yaml:"enabled"`
		Rate    string `yaml:"rate"`
		Burst   int    `yaml:"burst"`
	} `yaml:"rate_limit"`

	JWT struct {
		Enabled   bool   `yaml:"enabled"`
		SecretKey string `yaml:"secret_key"`
	} `yaml:"jwt"`

	Dashboard struct {
		Enabled bool   `yaml:"enabled"`
		Port    string `yaml:"port"`
	} `yaml:"dashboard"`

	REPL struct {
		Enabled bool `yaml:"enabled"`
	} `yaml:"repl"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	cfg := &Config{}
	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
