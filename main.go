package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func main() {
	proxy, err := NewProxy("https://api.github.com")
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/", ProxyRequestHandler(proxy))
	http.ListenAndServe(":5174", nil)
}

func NewProxy(targetHost string) (*httputil.ReverseProxy, error) {
	target, err := url.Parse(targetHost)
	if err != nil {
		return nil, err
	}

	proxy := httputil.NewSingleHostReverseProxy(target)

	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		modifyRequest(req, target)
	}

	proxy.ModifyResponse = modifyResponse()
	proxy.ErrorHandler = errorHandler()
	return proxy, nil
}

func ProxyRequestHandler(proxy *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		proxy.ServeHTTP(w, r)
	}
}

func modifyResponse() func(*http.Response) error {
	return func(resp *http.Response) error {
		fmt.Println("modifyResponse")
		return nil
	}
}

func modifyRequest(req *http.Request, target *url.URL) error {
	req.Host = target.Host

	req.Header.Set("Authorization", req.Header.Get("token"))
	return nil
}

func errorHandler() func(http.ResponseWriter, *http.Request, error) {
	return func(w http.ResponseWriter, req *http.Request, err error) {
		fmt.Printf("Got error while modifying response: %v \n", err)
		return
	}
}
