package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

var Config Configuration

type Configuration struct {
	Methods []MethodsConfig `json:"methods"`
}

type MethodsConfig struct {
	Name      string   `json:"name"`
	ProxyTo   []string `json:"proxy-to"`
	RateLimit int      `json:"rate-limit"`
}

func LoadConfigFile(filePath string) Configuration {
	file, errReadFile := ioutil.ReadFile(filePath)
	if errReadFile != nil {
		log.Panicf("Could not read the file (%v): %v", file, errReadFile)
	}
	return LoadConfig(string(file))
}

func LoadConfig(content string) Configuration {
	config := Configuration{}
	errJson := json.Unmarshal([]byte(content), &config)
	if errJson != nil {
		log.Panicf("Could not read json content: %v", errJson)
	}
	return config
}
