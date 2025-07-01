package core

import "net/http"

type Plugin interface {
	Name() string
	Init(config map[string]interface{}) error
	Execute(http.ResponseWriter, *http.Request) error
}
