package mist

import (
	"log"
	"net/http"

	"github.com/elos/ehttp/serve"
	"github.com/elos/mist/services"
	"github.com/gorilla/context"
)

type Middleware struct {
	Cors serve.Middleware

	Log serve.Middleware
}

type Services struct {
	services.DB

	services.Twilio

	services.Views
}

type Mist struct {
	router serve.Router
	*Middleware
	*Services
}

func New(m *Middleware, s *Services) *Mist {
	router := router(m, s)

	if m.Cors == nil {
		log.Fatal("Middleware Cors is nil")
	}

	if m.Log == nil {
		log.Fatal("Middleware Log is nil")
	}

	if s.DB == nil {
		log.Fatal("Service DB is nil")
	}

	if s.Twilio == nil {
		log.Fatal("Service Twilio is nil")
	}

	if s.Views == nil {
		log.Fatal("Service Views is nil")
	}

	return &Mist{
		router:     router,
		Middleware: m,
		Services:   s,
	}
}

func (mist *Mist) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	context.ClearHandler(http.HandlerFunc(mist.router.ServeHTTP)).ServeHTTP(w, r)
}
