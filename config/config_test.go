package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLoadConfig_Valid(t *testing.T) {
	json := `{
		"base-url": "/",
		"methods": [
			{
			  "name": "searchrawtransactions",
			  "proxy-to": [
				"http://localhost:2021",
				"http://localhost:2022"
			  ],
			  "rate-limit": 10
			},
			{
			  "name": ".*",
			  "proxy-to": [
				"http://localhost:2023"
			  ],
			  "rate-limit": 100
			}
		]
	}`
	config := LoadConfig(json)
	assert.Equal(t, config.BaseURL, "/")
	assert.Equal(t, config.Methods[0].Name, "searchrawtransactions")
	assert.Equal(t, config.Methods[0].ProxyTo[0], "http://localhost:2021")
	assert.Equal(t, config.Methods[0].ProxyTo[1], "http://localhost:2022")
	assert.Equal(t, config.Methods[0].RateLimit, 10)
	assert.Equal(t, config.Methods[1].Name, ".*")
	assert.Equal(t, config.Methods[1].ProxyTo[0], "http://localhost:2023")
	assert.Equal(t, config.Methods[1].RateLimit, 100)
	assert.Equal(t, len(config.Methods), 2)
	assert.Equal(t, len(config.Methods[0].ProxyTo), 2)
	assert.Equal(t, len(config.Methods[1].ProxyTo), 1)
}

func TestLoadConfig_MissingAttributes(t *testing.T) {
	json := `{
		"base-url": "/",
		"methods": [
			{
			  "name": "searchrawtransactions",
			  "proxy-to": [
				"http://localhost:2021",
				"http://localhost:2022"
			  ]
			},
			{
			  "name": ".*"
			}
		]
	}`
	config := LoadConfig(json)
	assert.Equal(t, config.BaseURL, "/")
	assert.Equal(t, config.Methods[0].Name, "searchrawtransactions")
	assert.Equal(t, config.Methods[0].ProxyTo[0], "http://localhost:2021")
	assert.Equal(t, config.Methods[0].ProxyTo[1], "http://localhost:2022")
	assert.Equal(t, config.Methods[0].RateLimit, 0)
	assert.Equal(t, config.Methods[1].Name, ".*")
	assert.Nil(t, config.Methods[1].ProxyTo)
	assert.Equal(t, config.Methods[1].RateLimit, 0)
	assert.Equal(t, len(config.Methods), 2)
	assert.Equal(t, len(config.Methods[0].ProxyTo), 2)
	assert.Equal(t, len(config.Methods[1].ProxyTo), 0)
}
