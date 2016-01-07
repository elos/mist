package services

import (
	"io"

	"github.com/elos/data"
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
