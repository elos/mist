package services

import (
	"io"

	"github.com/elos/data"
	"github.com/elos/echo"
	"github.com/elos/ehttp/templates"
	"github.com/subosito/twilio"
)

type DB interface {
	data.DB
}

type Twilio interface {
	Send(to, body string) (*twilio.Message, *twilio.Response, error)
}

type Views interface {
	Execute(w io.Writer, name templates.Name, data interface{}) error
}

type Texts interface {
	// Input is a send only channel on which this server can notify
	// that it has recieved a message
	Input() chan<- *echo.Message
}
