package routes

import "github.com/elos/ehttp/serve"

func TestGET(c *serve.Conn) {
	c.Write([]byte("hello"))
}
