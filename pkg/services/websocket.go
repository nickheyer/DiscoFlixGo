package services

import (
	"errors"

	"github.com/labstack/echo/v4"
	"github.com/nickheyer/DiscoFlixGo/pkg/log"
	"golang.org/x/net/websocket"
)

// WsClient manages WebSocket connections and broadcasts messages to all connections.
type WsClient struct {
	pool        *ConnectionPool
	echo        *echo.Echo
	addCh       chan *websocket.Conn
	removeCh    chan *websocket.Conn
	broadcastCh chan string
}

// ConnectionPool manages WebSocket connections.
type ConnectionPool struct {
	connections map[*websocket.Conn]struct{}
}

// NewWsClient creates a new WsClient and sets up the WebSocket route.
func NewWsClient(e *echo.Echo) *WsClient {
	client := &WsClient{
		pool: &ConnectionPool{
			connections: make(map[*websocket.Conn]struct{}),
		},
		echo:        e,
		addCh:       make(chan *websocket.Conn),
		removeCh:    make(chan *websocket.Conn),
		broadcastCh: make(chan string),
	}
	go client.manageConnections()
	return client
}

// manageConnections handles connection management and broadcasting messages.
func (client *WsClient) manageConnections() {
	for {
		select {
		case conn := <-client.addCh:
			client.pool.connections[conn] = struct{}{}
			log.Default().Info("WebSocket connection established.")
		case conn := <-client.removeCh:
			delete(client.pool.connections, conn)
			log.Default().Info("WebSocket connection removed.")
		case msg := <-client.broadcastCh:
			for conn := range client.pool.connections {
				go func(c *websocket.Conn) {
					if err := websocket.Message.Send(c, msg); err != nil {
						log.Default().Error("Unable to send message via WebSocket: ", "error", err)
					}
				}(conn)
			}
		}
	}
}

// AddConnectionToPool adds a WebSocket connection to the pool and listens for messages.
func (client *WsClient) AddConnectionToPool(ctx echo.Context) error {
	websocket.Handler(func(ws *websocket.Conn) {
		client.addCh <- ws
		defer func() { client.removeCh <- ws }()
		client.broadcastCh <- "client joined"
		client.listen(ws)
	}).ServeHTTP(ctx.Response(), ctx.Request())
	return nil
}

// listen monitors and relays WebSocket messages from a connection.
func (client *WsClient) listen(wsConn *websocket.Conn) error {
	msg := ""
	for {
		if err := websocket.Message.Receive(wsConn, &msg); err != nil {
			return errors.New("error while receiving WebSocket message")
		}
		client.broadcastCh <- msg
	}
}

// Broadcast sends a message to all connected client
func (client *WsClient) Broadcast(msg string) {
	client.broadcastCh <- msg
}

func (client *WsClient) Remove(ws *websocket.Conn) {
	client.removeCh <- ws
}

// GetEchoInstance returns the Echo instance associated with the WsClient.
func (client *WsClient) GetEchoInstance() *echo.Echo {
	return client.echo
}

func (client *WsClient) GetConnections() map[*websocket.Conn]struct{} {
	return client.Pool().connections
}

func (client *WsClient) Pool() *ConnectionPool {
	return client.pool
}
