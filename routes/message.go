package routes

import (
	"log"

	"github.com/elos/echo"
	"github.com/elos/ehttp/serve"
	"github.com/elos/mist/services"
)

func MessagePOST(c *serve.Conn, db services.DB, twilio services.Twilio) {
	m, err := echo.Extract(c)

	if err != nil {
		log.Fatal(err)
	}

	echo.Handle(m, twilio)
}

func MessageOPTIONS(c *serve.Conn) {
	c.WriteHeader(200)
}
