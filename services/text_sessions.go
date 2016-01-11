package services

import (
	"log"
	"strings"
	"time"

	"github.com/elos/echo"
	"github.com/elos/elos/command"
	"github.com/mitchellh/cli"
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
				// to prevent deadlock between a session that we are
				// trying to send to bailing
				go func() {
					ts.sessions[e.From] <- e.Body
				}()
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
				// Construct a new CLI with our name and version
				c := cli.NewCLI("elos [mist]", "0.0.1")

				// Tokenize the message
				c.Args = strings.Split(text, " ")

				ui, out := echo.NewTextUI(input, "TODO")
				go func() {
					for s := range out {
						twilio.Send(from, s)
					}
				}()

				// Initialize the commands
				c.Commands = map[string]cli.CommandFactory{
					"note": func() (cli.Command, error) {
						return &command.NoteCommand{
							Ui: ui,
							Config: &command.Config{
								DB: "localhost",
							},
						}, nil
					},
				}

				// now we block on run, so someone else can pop from input,
				// namely the UI
				_, err := c.Run()
				if err != nil {
					log.Print(err.Error())
				}
			case <-time.After(5 * time.Minute):
				twilio.Send(from, "Session time out")
				bail <- from
				break Run
			}
		}

	}()

	return input
}
