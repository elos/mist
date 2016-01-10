package services

import (
	"time"

	"github.com/elos/echo"
)

func NewTexts(twilio Twilio) *TextSessions {
	t := &TextSessions{
		input:    make(chan *echo.Message),
		sessions: make(map[string]chan<- string),
	}

	return t
}

type TextSessions struct {
	input    chan *echo.Message
	sessions map[string]chan<- string
}

func (ts *TextSessions) Input() chan<- *echo.Message {
	return ts.input
}

func (ts *TextSessions) Run(twilio Twilio) chan<- struct{} {
	done := make(chan struct{})
	bails := make(chan string)

	go func() {
	Run:
		for {
			select {
			case e := <-ts.input:
				_, ok := ts.sessions[e.From]
				if !ok {
					ts.sessions[e.From] = session(e.From, bails, twilio)
				}
				ts.sessions[e.From] <- e.Body
			case s := <-bails:
				if _, ok := ts.sessions[s]; ok {
					delete(ts.sessions, s)
				}
			case <-done:
				break Run
			}
		}
	}()

	return done
}

func session(from string, bail chan<- string, twilio Twilio) chan<- string {
	input := make(chan string)

	go func() {
	Run:
		for {
			select {
			case text := <-input:
				twilio.Send(from, text)

			case <-time.After(5 * time.Minute):
				if from != "" {
					twilio.Send(from, "Session time out")
				}
				bail <- from
				break Run
			}
		}
	}()

	return input
}
