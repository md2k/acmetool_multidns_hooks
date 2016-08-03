package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	yaml "gopkg.in/yaml.v2"
)

// FindDomainProvider() Domain Details Structure
type DomainDetails struct {
	Provider string
	Account  string
	AuthData map[string]string
}

// Yaml Config Structure
type HookConfig struct {
	Providers map[string]Provider `yaml:"providers"`
}

type Provider struct {
	Accounts map[string]Account `yaml:"accounts"`
}

type Account struct {
	AuthData map[string]string `yaml:"authdata"`
	Domains  []string          `yaml:"domains"`
}

func InitConfig(acme_path string) (cfg *HookConfig, err error) {
	// Get My Name (AppName)
	_, execName := filepath.Split(os.Args[0])
	// Prepare Config full-path location
	cfg_file := fmt.Sprintf("%s/conf/%s.yaml", acme_path, execName)
	// Read Yaml Config
	yamlFile, err := ioutil.ReadFile(cfg_file)
	if err != nil {
		return nil, fmt.Errorf("%s", err)
	}

	// Parse Yaml Config
	// var cfg HookConfig
	err = yaml.Unmarshal(yamlFile, &cfg)
	if err != nil {
		return nil, fmt.Errorf("%s", err)
	}

	return cfg, nil
}
