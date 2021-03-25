package routes

import (
	"github.com/elos/data"
	"github.com/elos/ehttp/serve"
	"github.com/elos/mist/services"
	"golang.org/x/net/websocket"
)

func ws(db data.DB) func(ws *websocket.Conn) {
	return func(ws *websocket.Conn) {
		services.WebSocketSession(db, ws)
	}
}

func WsGET(c *serve.Conn, db data.DB) {
	websocket.Handler(ws(db)).ServeHTTP(c.ResponseWriter(), c.Request())
}

func WsOPTIONS(c *serve.Conn) {
}
