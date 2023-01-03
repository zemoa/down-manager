package service

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/olahol/melody"
)

type WebSocketService struct {
	m *melody.Melody
}

func NewWebSocket() *WebSocketService {
	wss := new(WebSocketService)
	wss.m = melody.New()
	wss.m.HandleConnect(func(s *melody.Session) {
		id := uuid.NewString()
		log.Printf("Client <%s> connected", id)
		s.Set("id", id)
	})
	wss.m.HandleDisconnect(func(s *melody.Session) {
		id, _ := s.Get("id")
		log.Printf("client <%s> disconnected", id)
	})
	return wss
}

func (wss *WebSocketService) WebSocket() func(c *gin.Context) {
	return func(c *gin.Context) {
		wss.m.HandleRequest(c.Writer, c.Request)
	}
}

func (wss *WebSocketService) BroadcastMessage(msg []byte) {
	wss.m.Broadcast(msg)
}
