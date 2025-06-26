package core

import "net/http"

type RequestContext struct {
    Writer  http.ResponseWriter
    Request *http.Request
    Params  map[string]string
}
