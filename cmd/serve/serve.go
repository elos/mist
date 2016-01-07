package main

import (
	"flag"
	"log"

	"github.com/elos/autonomous"
	emiddleware "github.com/elos/ehttp/middleware"
	"github.com/elos/ehttp/serve"
	"github.com/elos/mist"
	mistmiddleware "github.com/elos/mist/middleware"
	"github.com/elos/models"
	"github.com/subosito/twilio"
)

var (
	AccountSid = "AC76d4c9975dfb641d9ae711c2f795c5a2"
	AuthToken  = "9ab82f10b0b6187d2c71589c46c96da6"
	From       = "+16503810349"
)

type TwilioService struct {
	client *twilio.Client
}

func (ts *TwilioService) Send(to, body string) (*twilio.Message, *twilio.Response, error) {
	return ts.client.Messages.SendSMS(From, to, body)
}

func main() {
	var (
		mongo = flag.String("mongo", "localhost", "Address of mongo instance")
	)

	flag.Parse()

	hub := autonomous.NewHub()
	go hub.Start()
	hub.WaitStart()

	db, err := models.MongoDB(*mongo)
	if err != nil {
		log.Fatal(err)
	}
	// Setup Middleware
	middleware := &mist.Middleware{
		Cors: mistmiddleware.NewCors(),
		Log:  emiddleware.LogRequest,
	}

	// Initialize twilio client
	c := twilio.NewClient(AccountSid, AuthToken, nil)

	// Setup Services
	services := &mist.Services{
		DB:     db,
		Twilio: &TwilioService{client: c},
	}

	mist := mist.New(middleware, services)

	serveOptions := &serve.Opts{
		Handler: mist,
	}

	server := serve.New(serveOptions)
	hub.StartAgent(server)

	go autonomous.HandleIntercept(hub.Stop)
	hub.WaitStop()
}
