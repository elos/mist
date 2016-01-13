package services

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/elos/data"
	"github.com/elos/echo"
	"github.com/elos/elos/command"
	"github.com/elos/models"
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

func (ts *TextSessions) Run(db data.DB, twilio Twilio) chan<- struct{} {
	done := make(chan struct{})
	bails := make(chan string)

	go func() {
	Run:
		for {
			select {
			case e := <-ts.input:
				_, ok := ts.sessions[e.From]
				if !ok {
					ts.sessions[e.From] = session(e.From, bails, db, twilio)
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

func session(from string, bail chan<- string, db data.DB, twilio Twilio) chan<- string {
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

				ui, out := echo.NewTextUI(input)
				go func() {
					for s := range out {
						twilio.Send(from, s)
					}
				}()
				var p *models.Person
				var userID string
				match := false

				q := db.NewQuery(models.PersonKind)
				q.Select(data.AttrMap{
					"phone": from,
				})
				iter, err := q.Execute()
				if err != nil {
					log.Print(err)
					goto Bail
				}

				p = models.NewPerson()

				for iter.Next(p) {
					match = true
				}

				if err := iter.Close(); err != nil {
					log.Print(err)
					goto Bail
				}

				if !match {
					ui.Output("Looks like you are new to elos")
					ui.Output("Welcome!")
					ui.Output("Let's get you an account")

					u := models.NewUser()
					u.SetID(db.NewID())
					u.CreatedAt = time.Now()
					u.UpdatedAt = time.Now()
					if err := db.Save(u); err != nil {
						log.Print(err)
						goto Bail
					}

					p.SetID(db.NewID())
					p.CreatedAt = time.Now()
					p.UpdatedAt = time.Now()
					p.Phone = from
					p.OwnerId = u.Id

					if err = db.Save(p); err != nil {
						log.Print(err)
						goto Bail
					}

					ui.Output(fmt.Sprintf("You're ID is: %s", u.ID()))
					userID = u.ID().String()
				} else {
					userID = p.OwnerId
				}

				// Initialize the commands
				c.Commands = map[string]cli.CommandFactory{
					"note": func() (cli.Command, error) {
						return &command.NoteCommand{
							Ui: ui,
							Config: &command.Config{
								DB:     "localhost",
								UserID: userID,
							},
							DB: db,
						}, nil
					},
					"todo": func() (cli.Command, error) {
						return &command.TodoCommand{
							UI:     ui,
							UserID: userID,
							DB:     db,
						}, nil
					},
				}

				// now we block on run, so someone else can pop from input,
				// namely the UI
				_, err = c.Run()
				if err != nil {
					log.Print(err.Error())
				}
				break

			Bail:
				bail <- from
				break Run
			case <-time.After(5 * time.Minute):
				twilio.Send(from, "Session time out")
				bail <- from
				break Run
			}
		}

	}()

	return input
}
