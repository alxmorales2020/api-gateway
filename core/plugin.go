package core

type Plugin interface {
    Init(config map[string]interface{}) error
    Handle(ctx *RequestContext) error
    Name() string
}
