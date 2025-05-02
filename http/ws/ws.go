package ws

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

// RegisterWebsocket sets up the WebSocket endpoint.
func RegisterWebsocket(e *echo.Echo) {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}

	e.GET("/ws", func(c echo.Context) error {
		conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
		if err != nil {
			return err
		}
		defer conn.Close()

		conn.WriteMessage(websocket.TextMessage, []byte("connected"))
		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				break
			}
			conn.WriteMessage(websocket.TextMessage, msg)
		}
		return nil
	})
}
