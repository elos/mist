package routes

import (
	"github.com/elos/ehttp/serve"
	"github.com/elos/mist/services"
)

func RegisterGET(c *serve.Conn) {
	c.Write([]byte("Register Get"))
}

func RegisterPOST(c *serve.Conn, db services.DB) {
	c.Write([]byte("Register Post"))
}
