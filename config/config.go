package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

// Config var with the configuration parsed from the .json file
var Config Configuration

// Configuration type representing the .json config file
type Configuration struct {
	BaseURL   string          `json:"base-url"`
	SSL       bool            `json:"ssl,omitempty"`
	SSLDomain string          `json:"ssl-domain,omitempty"`
	Methods   []MethodsConfig `json:"methods"`
}

// MethodsConfig type representing the config for each JSON-RPC method
type MethodsConfig struct {
	Name      string   `json:"name"`
	ProxyTo   []string `json:"proxy-to"`
	RateLimit int      `json:"rate-limit"`
}

// LoadConfigFile loads the config file into an instance of Configuration type
func LoadConfigFile(filePath string) Configuration {
	file, errReadFile := ioutil.ReadFile(filePath)
	if errReadFile != nil {
		log.Panicf("Could not read the file (%v): %v", file, errReadFile)
	}
	return LoadConfig(string(file))
}

// LoadConfig loads the content (as JSON) into an instance of Configuration type
func LoadConfig(content string) Configuration {
	config := Configuration{}
	errJSON := json.Unmarshal([]byte(content), &config)
	if errJSON != nil {
		log.Panicf("Could not read json content: %v", errJSON)
	}
	return config
}
