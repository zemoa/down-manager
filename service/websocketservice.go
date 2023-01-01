package service

import (
	"log"
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/google/uuid"
)

type WebSocketService struct {
	activeConnections map[uuid.UUID]ConnWrapper
}

func (wss *WebSocketService) Init() {
	wss.activeConnections = make(map[uuid.UUID]ConnWrapper)
}

func (wss *WebSocketService) WebSocket() func(c *gin.Context) {
	return func(c *gin.Context) {
		conn, _, _, err := ws.UpgradeHTTP(c.Request, c.Writer)
		if err != nil {
			c.Writer.WriteHeader(http.StatusInternalServerError)
		} else {
			uuid, _ := uuid.NewUUID()
			var connWrapper = ConnWrapper{id: uuid, conn: conn}
			wss.activeConnections[connWrapper.id] = connWrapper
			connWrapper.init()
			go func() {
				defer wss.closeChannel(&connWrapper)

				for {
					if !connWrapper.listen() {
						return
					}
				}
			}()
		}
	}
}

func (wss *WebSocketService) closeChannel(connWrapper *ConnWrapper) {
	connWrapper.close()
	delete(wss.activeConnections, connWrapper.id)
}

type ConnWrapper struct {
	id   uuid.UUID
	conn net.Conn
}

func (connWrapper *ConnWrapper) init() {
	log.Printf("New client connected through websocket")
}

func (connWrapper *ConnWrapper) listen() (isOpen bool) {
	_, op, err := wsutil.ReadClientData(connWrapper.conn)
	if op == ws.OpClose || err != nil {
		log.Printf("Client <%s> quit %s", connWrapper.id, err)
		return false
	}
	return true
}

func (connWrapper *ConnWrapper) Write() {
	connWrapper.conn.Close()
}

func (connWrapper *ConnWrapper) close() {
	log.Printf("Client <%s> quit", connWrapper.id)
	connWrapper.conn.Close()
}
