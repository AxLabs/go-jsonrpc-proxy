package main

import (
	"github.com/AxLabs/go-jsonrpc-proxy/config"
	"github.com/AxLabs/go-jsonrpc-proxy/server"
	"math/rand"
	"net/http"
	"os"
	"time"
)

func main() {
	initRand()

	configFilePath := getEnv("JSONRPC_PROXY_CONFIG_FILE", "jsonrpc-proxy.json")
	config := config.LoadConfigFile(configFilePath)
	server.LoadMap(config)

	http.HandleFunc(config.BaseURL, server.HandleRequestAndRedirect)
	if err := http.ListenAndServe(getListenAddress(), nil); err != nil {
		panic(err)
	}
}

// Get env var or default
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// Get the port to listen on
func getListenAddress() string {
	port := getEnv("PORT", "2000")
	return ":" + port
}

func initRand() {
	// initialize global pseudo random generator
	rand.Seed(time.Now().Unix())
}
