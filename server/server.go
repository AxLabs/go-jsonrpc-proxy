package server

import (
	"bytes"
	"github.com/AxLabs/go-jsonrpc-proxy/config"
	"github.com/patrickmn/go-cache"
	"github.com/teambition/jsonrpc-go"
	"gopkg.in/square/go-jose.v2/json"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"net/url"
	"regexp"
	"time"
)

var methodNameCache *cache.Cache

var methodPriorityListOrder []MethodRegExp

type MethodRegExp struct {
	Name       string
	NameRegexp regexp.Regexp
	ProxyTo    []string
}

func NewMethodRegExp(name string, nameRegexp regexp.Regexp, proxyTo []string) MethodRegExp {
	return MethodRegExp{
		Name:       name,
		NameRegexp: nameRegexp,
		ProxyTo:    proxyTo,
	}
}

func LoadMap(config config.Configuration) {
	methodNameCache = cache.New(5*time.Minute, 5*time.Minute)
	methodPriorityListOrder = []MethodRegExp{}
	for _, method := range config.Methods {
		compiledMethodName, errCompile := regexp.Compile(method.Name)
		if errCompile != nil {
			log.Panicf("Config file contains an invalid regex in method's name: %v", method.Name)
		}
		methodRegExpObj := NewMethodRegExp(method.Name, *compiledMethodName, method.ProxyTo)
		methodPriorityListOrder = append(methodPriorityListOrder, methodRegExpObj)
	}
}

// Serve a reverse proxy for a given url
func serveReverseProxy(target string, res http.ResponseWriter, req *http.Request) {
	// parse the url
	url, _ := url.Parse(target)

	// create the reverse proxy
	proxy := httputil.NewSingleHostReverseProxy(url)

	// Update the headers to allow for SSL redirection
	req.URL.Host = url.Host
	req.URL.Scheme = url.Scheme
	req.Header.Set("X-Forwarded-Host", req.Header.Get("Host"))
	req.Host = url.Host

	// The ServeHttp is non-blocking and uses a go routine under the hood
	proxy.ServeHTTP(res, req)
}

func requestBody(request *http.Request) *bytes.Buffer {
	// Read body to buffer
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		log.Printf("Error reading body: %v", err)
		panic(err)
	}

	request.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	return bytes.NewBuffer(body)
}

func parseRequestBody(request *http.Request) *jsonrpc.RPC {
	buffered := requestBody(request)
	req, _ := jsonrpc.Parse(buffered.Bytes())
	if req == nil {
		log.Printf("could not parse the JSON-RPC")
	}
	return req
}

// HandleRequestAndRedirect given a request send it to the appropriate url
func HandleRequestAndRedirect(res http.ResponseWriter, req *http.Request) {
	requestPayload := parseRequestBody(req)
	url, errRedir := getRedirectTo(requestPayload)

	if errRedir != nil {
		errObjBytes, errJsonMarshal := json.Marshal(errRedir)
		if errJsonMarshal != nil {
			internalError := jsonrpc.InternalError()
			internalErrorBytes, _ := json.Marshal(internalError)
			res.Write(internalErrorBytes)
			return
		}
		res.Write(errObjBytes)
		return
	}

	serveReverseProxy(*url, res, req)
}

func getRedirectTo(req *jsonrpc.RPC) (*string, *jsonrpc.ErrorObj) {
	if methodNameCache == nil {
		return nil, jsonrpc.InternalError("cache not loaded")
	}
	if value, ok := methodNameCache.Get(req.Method); ok {
		methodRegEx := value.(MethodRegExp)
		randomRedirectTo := getRandomElem(methodRegEx.ProxyTo)
		return &randomRedirectTo, nil
	}
	for _, method := range methodPriorityListOrder {
		if method.NameRegexp.MatchString(req.Method) {
			methodNameCache.SetDefault(req.Method, method)
			randomRedirectTo := getRandomElem(method.ProxyTo)
			return &randomRedirectTo, nil
		}
	}
	return nil, jsonrpc.MethodNotFound()
}

func getRandomElem(array []string) string {
	return array[rand.Intn(len(array))]
}
