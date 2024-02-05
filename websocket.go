package main

import (
	"fmt"
	"io"

	"golang.org/x/net/websocket"
)

// WebSocket server used for hot reload of the front-end

type WebSocketServer struct {
	conns map[*websocket.Conn]bool
}

func NewWebSocketServer() *WebSocketServer {
	return &WebSocketServer{conns: make(map[*websocket.Conn]bool)}
}

func (s *WebSocketServer) handleWS(ws *websocket.Conn) {
	s.conns[ws] = true

	s.readLoop(ws)
}

func (s *WebSocketServer) readLoop(ws *websocket.Conn) {
	buf := make([]byte, 1024)
	for {
		n, err := ws.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("Error reading message", err)
			break
		}
		msg := string(buf[:n])
		fmt.Println(msg)
	}
}
