package main

import (
	"github.com/AxLabs/go-jsonrpc-proxy/config"
	"github.com/AxLabs/go-jsonrpc-proxy/server"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/shibukawa/configdir"
	"golang.org/x/crypto/acme/autocert"
	"math/rand"
	"net/http"
	"os"
	"time"
)

func main() {
	// initialize random
	initRand()

	// load config
	configFilePath := getEnvVar("JSONRPC_PROXY_CONFIG_FILE", "jsonrpc-proxy.json")
	config := config.LoadConfigFile(configFilePath)
	server.LoadMap(config)

	sslEnabled := getEnvVar("JSONRPC_PROXY_SSL", "off")
	sslDomain := getEnvVar("JSONRPC_PROXY_SSL_DOMAIN", "")
	httpRedirectToHTTPS := getEnvVar("JSONRPC_PROXY_REDIRECT_TO_HTTPS", "on")
	configDir := getConfigDir()

	e := echo.New()

	if sslEnabled == "on" && httpRedirectToHTTPS == "on" {
		e.Pre(middleware.HTTPSRedirect())
	}

	e.AutoTLSManager.Cache = autocert.DirCache(configDir)
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())

	e.POST("/", echo.WrapHandler(http.HandlerFunc(server.HandleRequestAndRedirect)))

	if sslEnabled == "on" {
		if len(sslDomain) == 0 {
			e.Logger.Fatalf("You should specify the domain name for the SSL.")
		}
		e.AutoTLSManager.HostPolicy = autocert.HostWhitelist(sslDomain)
		e.Logger.Fatal(e.StartAutoTLS(getListenAddress()))
	} else {
		e.Logger.Fatal(e.Start(getListenAddress()))
	}
}

// getEnvVar get an environment variable value, or returns a fallback value
func getEnvVar(key, fallback string) string {
	value, _ := getEnvVarAndIfExists(key, fallback)
	return value
}

// getEnvVarAndIfExists Retrieves an environment variable value, or returns a default (fallback) value
// It also returns true or false if the env variable exists or not
func getEnvVarAndIfExists(key, fallback string) (string, bool) {
	value, exists := os.LookupEnv(key)
	if len(value) == 0 {
		return fallback, exists
	}
	return value, exists
}

// getListenAddress get the port to listen on
func getListenAddress() string {
	port := getEnvVar("PORT", "2000")
	return ":" + port
}

// getConfigDir get the configuration dir
func getConfigDir() string {
	configDirs := configdir.New("axlabs", "go-jsonrpc-proxy")
	folders := configDirs.QueryFolders(configdir.Global)
	err := folders[0].MkdirAll()
	if err != nil {
		log.Panicf("Not possible to create folder %v: %v", folders[0].Path, err)
	}

	configPath, ok := getEnvVarAndIfExists("JSONRPC_PROXY_CONFIG_PATH", folders[0].Path)
	if ok {
		configDirs.LocalPath = configPath
		folders = configDirs.QueryFolders(configdir.Local)
	}
	err = folders[0].MkdirAll()
	if err != nil {
		log.Panicf("Not possible to create folder %v: %v", folders[0].Path, err)
	}
	return folders[0].Path
}

func initRand() {
	// initialize global pseudo random generator
	rand.Seed(time.Now().Unix())
}
