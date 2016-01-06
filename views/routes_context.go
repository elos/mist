package views

import (
	"fmt"

	"github.com/elos/ehttp/templates"
)

type RoutesContext struct {
}

func (r *RoutesContext) Message() string {
	return fmt.Sprintf("/message")
}

var routesContext = &RoutesContext{}

type context struct {
	Routes *RoutesContext
	Data   interface{}
}

func (c *context) WithData(d interface{}) templates.Context {
	return &context{
		Routes: c.Routes,
		Data:   d,
	}
}

var globalContext = &context{
	Routes: routesContext,
}
