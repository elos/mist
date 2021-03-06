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

func (r *RoutesContext) Register() string {
	return fmt.Sprintf("/register")
}

func (r *RoutesContext) Test() string {
	return fmt.Sprintf("/test")
}

func (r *RoutesContext) Ws() string {
	return fmt.Sprintf("/ws")
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
