package mist

import (
	"github.com/elos/ehttp/builtin"
	"github.com/elos/ehttp/serve"
	"github.com/elos/mist/routes"
)

func router(m *Middleware, s *Services) serve.Router {
	router := builtin.NewRouter()

	router.GET(routes.Test, func(c *serve.Conn) {

		if ok := m.Log.Inbound(c); !ok {
			return
		}

		if ok := m.Cors.Inbound(c); !ok {
			return
		}

		routes.TestGET(c)

		if ok := m.Cors.Outbound(c); !ok {
			return
		}

		if ok := m.Log.Outbound(c); !ok {
			return
		}

	})

	router.GET(routes.Register, func(c *serve.Conn) {

		if ok := m.Log.Inbound(c); !ok {
			return
		}

		if ok := m.Cors.Inbound(c); !ok {
			return
		}

		routes.RegisterGET(c, s.Views)

		if ok := m.Cors.Outbound(c); !ok {
			return
		}

		if ok := m.Log.Outbound(c); !ok {
			return
		}

	})

	router.POST(routes.Register, func(c *serve.Conn) {

		if ok := m.Log.Inbound(c); !ok {
			return
		}

		if ok := m.Cors.Inbound(c); !ok {
			return
		}

		routes.RegisterPOST(c, s.DB)

		if ok := m.Cors.Outbound(c); !ok {
			return
		}

		if ok := m.Log.Outbound(c); !ok {
			return
		}

	})

	router.POST(routes.Message, func(c *serve.Conn) {

		if ok := m.Log.Inbound(c); !ok {
			return
		}

		if ok := m.Cors.Inbound(c); !ok {
			return
		}

		routes.MessagePOST(c, s.DB, s.Twilio, s.Texts)

		if ok := m.Cors.Outbound(c); !ok {
			return
		}

		if ok := m.Log.Outbound(c); !ok {
			return
		}

	})

	router.OPTIONS(routes.Message, func(c *serve.Conn) {

		if ok := m.Log.Inbound(c); !ok {
			return
		}

		if ok := m.Cors.Inbound(c); !ok {
			return
		}

		routes.MessageOPTIONS(c)

		if ok := m.Cors.Outbound(c); !ok {
			return
		}

		if ok := m.Log.Outbound(c); !ok {
			return
		}

	})

	router.GET(routes.Ws, func(c *serve.Conn) {

		if ok := m.Log.Inbound(c); !ok {
			return
		}

		if ok := m.Cors.Inbound(c); !ok {
			return
		}

		routes.WsGET(c, s.DB)

		if ok := m.Cors.Outbound(c); !ok {
			return
		}

		if ok := m.Log.Outbound(c); !ok {
			return
		}

	})

	router.OPTIONS(routes.Ws, func(c *serve.Conn) {

		if ok := m.Log.Inbound(c); !ok {
			return
		}

		if ok := m.Cors.Inbound(c); !ok {
			return
		}

		routes.WsOPTIONS(c)

		if ok := m.Cors.Outbound(c); !ok {
			return
		}

		if ok := m.Log.Outbound(c); !ok {
			return
		}

	})

	return router
}
