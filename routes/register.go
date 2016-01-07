package routes

import (
	"github.com/elos/ehttp/serve"
	"github.com/elos/ehttp/templates"
	"github.com/elos/mist/services"
	"github.com/elos/mist/views"
)

func RegisterGET(c *serve.Conn, engine services.Views) {
	templates.CatchError(c, engine.Execute(c, views.Register, nil))
}

func RegisterPOST(c *serve.Conn, db services.DB) {
	c.Write([]byte("Register Post"))
}
