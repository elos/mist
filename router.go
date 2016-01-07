package mist

import (
	"path/filepath"
	"runtime"

	"github.com/elos/ehttp/builtin"
	"github.com/elos/ehttp/serve"
	"github.com/elos/mist/routes"
)

var root string

func init() {
	_, filename, _, _ := runtime.Caller(1)
	root = filepath.Dir(filename)
}

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

		routes.RegisterGET(c)

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

		routes.MessagePOST(c, s.DB, s.Twilio)

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

	return router
}
