package handlers

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/nickheyer/DiscoFlixGo/pkg/services"
)

const (
	routeNameWebsocketInit = "websocketConnect"
)

type (
	WsConn struct {
		websocket *services.WsClient
	}
)

func init() {
	Register(new(WsConn))
}

func (h *WsConn) Init(c *services.Container) error {
	h.websocket = c.Websocket
	return nil
}

func (h *WsConn) Routes(g *echo.Group) {
	g.GET("/ws", h.connect).Name = routeNameWebsocketInit
}

func (h *WsConn) connect(ctx echo.Context) error {
	fmt.Printf("%+v\n", h.websocket)
	h.websocket.addConnectionToPool(ctx)
	return nil
}
