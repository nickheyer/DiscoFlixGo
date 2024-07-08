package services

import (
	"errors"
	"sync"

	"github.com/labstack/echo/v4"
	"github.com/nickheyer/DiscoFlixGo/pkg/log"
	"golang.org/x/net/websocket"
)

// WsClient manages WebSocket connections and broadcasts messages to all connections.
type WsClient struct {
	pool *ConnectionPool
	echo *echo.Echo
}

// ConnectionPool manages WebSocket connections.
type ConnectionPool struct {
	sync.RWMutex
	connections map[*websocket.Conn]struct{}
}

// NewWsClient creates a new WsClient and sets up the WebSocket route.
func NewWsClient(e *echo.Echo) *WsClient {
	client := &WsClient{
		pool: &ConnectionPool{
			connections: make(map[*websocket.Conn]struct{}),
		},
		echo: e,
	}
	return client
}

// addConnectionToPool adds a WebSocket connection to the pool and listens for messages.
func (client *WsClient) AddConnectionToPool(ctx echo.Context) error {
	websocket.Handler(func(ws *websocket.Conn) {
		// Add connection to pool.
		client.pool.Lock()
		client.pool.connections[ws] = struct{}{}
		log.Ctx(ctx).Info("Websocket connection established with session.")

		// Defer removal until client disconnects
		defer func(connection *websocket.Conn) {
			client.pool.Lock()
			delete(client.pool.connections, connection)
			log.Ctx(ctx).Info("Removing Websocket connection from connection pool.")
			client.pool.Unlock()
		}(ws)

		client.pool.Unlock()
		client.Broadcast("client joined")
		client.listen(ws)
	}).ServeHTTP(ctx.Response(), ctx.Request())
	return nil
}

// listen monitors and relays ws messages from a connection.
func (client *WsClient) listen(wsConn *websocket.Conn) error {
	msg := ""
	for {
		if err := websocket.Message.Receive(wsConn, &msg); err != nil {
			return errors.New("error while receiving WebSocket msg")
		}
		if err := client.Broadcast(msg); err != nil {
			return errors.New("error while sending WebSocket msg")
		}
	}
}

func (client *WsClient) Broadcast(msg string) error {
	client.pool.RLock()
	defer client.pool.RUnlock()
	log.Default().Info("Broadcasting message to all clients: ", msg, "...")
	for connection := range client.pool.connections {
		if err := websocket.Message.Send(connection, msg); err != nil {
			return errors.New("unable to send message via websocket")
		}
	}
	return nil
}

func (client *WsClient) GetEchoInstance() *echo.Echo {
	return client.echo
}
