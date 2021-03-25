package services

import (
	"log"
	"strings"
	"time"

	"github.com/elos/data"
	"github.com/elos/echo"
	"github.com/elos/ehttp/sock"
	"github.com/elos/elos/command"
	"github.com/elos/models/access"
	"github.com/mitchellh/cli"
	"golang.org/x/net/context"
	"golang.org/x/net/websocket"
)

type CommandSocks struct {
}

func (cs *CommandSocks) Dispatch(ctx context.Context, c sock.Conn) {
}

type message struct {
	kind string `json:"kind"`
	body string `json:"body"`
}

func read(ws *websocket.Conn) <-chan *message {
	out := make(chan *message)

	go func() {

	Read:
		for {
			m := new(message)
			err := websocket.JSON.Receive(ws, m)
			if err != nil {
				break Read
			}

			out <- m
		}

		close(out)
	}()

	return out
}

func digest(c <-chan *message) <-chan string {
	out := make(chan string)

	go func() {
		for m := range c {
			out <- m.body
		}

		close(out)
	}()

	return out
}

func WebSocketSession(db data.DB, ws *websocket.Conn) {
	messages := read(ws)
	input := digest(messages)

	ui, out := echo.NewTextUI(input)
	go func() {
		for s := range out {
			log.Printf("Sending: '%s'", s)
			websocket.JSON.Send(ws, &message{body: s})
		}
	}()

	public, err := ui.Ask("Username")

	if err != nil {
		log.Print("Error asking for username: %s", err)
		return
	}

	private, err := ui.Ask("password")
	if err != nil {
		log.Print("Error asking for password: %s", err)
		return
	}

	log.Printf("(%s, %s)", public, private)

	cred, err := access.Authenticate(db, public, private)
	if err != nil {
		log.Printf("Error on authenticating %s", err)
		return
	}
	log.Printf("Credential: %+v", cred)

	u, err := cred.Owner(db)
	if err != nil {
		log.Printf("Error on finding credential's owner: %s", err)
		return
	}
	log.Printf("User: %+v", u)

	userID := u.Id

Run:
	for {
		select {
		case text := <-input:
			// Construct a new CLI with our name and version
			instance := cli.NewCLI("elos [mist]", "0.0.1")

			// Tokenize the message
			instance.Args = strings.Split(text, " ")

			// Initialize the commands
			instance.Commands = map[string]cli.CommandFactory{
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
			_, err := instance.Run()
			if err != nil {
				log.Print(err.Error())
			}
		case <-time.After(5 * time.Minute):
			break Run
		}
	}
	log.Print("done")
}
