package main

import (
	"context"
	"crypto/tls"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gorilla/websocket"
)

func main() {
	var (
		ctx    = context.Background()
		logger = g.Log()
	)

	// Configure WebSocket dialer with TLS
	dialer := websocket.Dialer{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	// Connect to WebSocket server using wss://
	ws, _, err := dialer.Dial("wss://127.0.0.1:8000/ws", nil)
	if err != nil {
		logger.Fatalf(ctx, "dial failed: %+v", err)
	}
	defer ws.Close()

	err = ws.WriteMessage(websocket.TextMessage, []byte("hello"))
	if err != nil {
		logger.Fatalf(ctx, "ws.WriteMessage failed: %+v", err)
	}
	_, msg, err := ws.ReadMessage()
	if err != nil {
		logger.Fatalf(ctx, "ws.ReadMessage failed: %+v", err)
		return
	}
	logger.Infof(ctx, `received message: %s`, msg)

	// Cleanly close the connection by sending a close message
	err = ws.WriteMessage(websocket.CloseMessage, []byte("going to close"))
	if err != nil {
		logger.Fatalf(ctx, "ws.WriteMessage failed: %+v", err)
	}
}
