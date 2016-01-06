package services

import (
	"github.com/elos/data"
	"github.com/subosito/twilio"
)

type DB interface {
	data.DB
}

type Twilio interface {
	Send(to, body string) (*twilio.Message, *twilio.Response, error)
}
