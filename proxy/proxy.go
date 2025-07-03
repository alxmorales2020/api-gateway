package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

func NewReverseProxy(target string) (*httputil.ReverseProxy, error) {
	// Parse the target URL
	targetURL, err := url.Parse(target)
	if err != nil {
		return nil, err
	}

	// Create a new reverse proxy
	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	// Modify the request before sending it to the target
	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		// Call the original director to set the target URL
		originalDirector(req)

		req.Host = targetURL.Host
	}

	proxy.ErrorHandler = func(writer http.ResponseWriter, request *http.Request, err error) {
		http.Error(writer, "Upstream error: "+err.Error(), http.StatusBadGateway)
	}
	return proxy, nil
}
