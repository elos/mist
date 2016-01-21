package main

import (
	"flag"
	"log"
	"os"

	"github.com/elos/autonomous"
	emiddleware "github.com/elos/ehttp/middleware"
	"github.com/elos/ehttp/serve"
	"github.com/elos/ehttp/templates"
	"github.com/elos/gaia"
	gservices "github.com/elos/gaia/services"
	"github.com/elos/mist"
	mistmiddleware "github.com/elos/mist/middleware"
	"github.com/elos/mist/services"
	"github.com/elos/mist/views"
	"github.com/elos/models"
	"github.com/subosito/twilio"
	"golang.org/x/net/context"
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
	twilio := &TwilioService{client: c}

	sessions := services.NewTexts(twilio)
	killSessions := sessions.Run(db, twilio)

	// Setup Services
	services := &mist.Services{
		DB:     db,
		Twilio: twilio,
		Views:  templates.NewEngine(views.TemplatesDir, &views.TemplatesSet),
		Texts:  sessions,
	}

	services.Views.(*templates.Engine).Parse()

	mist := mist.New(middleware, services)

	serveOptions := &serve.Opts{
		Handler: mist,
	}

	server := serve.New(serveOptions)
	hub.StartAgent(server)

	mux := gservices.NewSMSMux()
	go mux.Start(context.TODO(), db, gservices.SMSFromTwilio(c, From))
	gaia := gaia.New(new(gaia.Middleware), &gaia.Services{
		SMSCommandSessions: mux,
		DB:                 db,
		Logger:             gservices.NewLogger(os.Stderr),
	})

	gaiaServeOptions := &serve.Opts{
		Port:    8080,
		Handler: gaia,
	}

	gaiaServer := serve.New(gaiaServeOptions)
	hub.StartAgent(gaiaServer)

	go autonomous.HandleIntercept(hub.Stop)
	hub.WaitStop()
	killSessions <- struct{}{}
}
