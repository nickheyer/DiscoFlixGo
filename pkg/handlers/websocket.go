package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/nickheyer/DiscoFlixGo/pkg/services"
)

const (
	routeNameWebsocketInit = "websocketConnect"
)

type (
	WsConn struct {
		ws *services.WsClient
	}
)

func init() {
	Register(new(WsConn))
}

func (h *WsConn) Init(c *services.Container) error {
	h.ws = c.Ws
	return nil
}

func (h *WsConn) Routes(g *echo.Group) {
	g.GET("/ws", h.connect).Name = routeNameWebsocketInit
}

func (h *WsConn) connect(ctx echo.Context) error {
	h.ws.AddConnectionToPool(ctx)
	return nil
}
