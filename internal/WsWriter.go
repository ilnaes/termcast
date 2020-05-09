package internal

import (
	"errors"

	"github.com/gorilla/websocket"
)

// an io.Writer wrapper on
// a websocket connection
type WsWriter struct {
	Conn *websocket.Conn
}

func (ww WsWriter) Write(p []byte) (n int, err error) {
	if ww.Conn == nil {
		return 0, errors.New("No socket connection")
	}
	err = ww.Conn.WriteMessage(websocket.TextMessage, p)
	if err != nil {
		return 0, err
	}

	return len(p), nil
}
