package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

// TimeoutMiddleware is used to handle request timeouts
func TimeoutMiddleware(httpTimeout time.Duration, wsTimeout time.Duration) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			var timeout time.Duration

			// Check if the request is a WebSocket request
			if isWebSocketRequest(c.Request()) {
				timeout = wsTimeout
			} else {
				timeout = httpTimeout
			}

			ctx, cancel := context.WithTimeout(c.Request().Context(), timeout)
			defer cancel()

			// Create a new request with the timeout context
			req := c.Request().WithContext(ctx)
			c.SetRequest(req)

			done := make(chan error, 1)
			go func() {
				done <- next(c)
			}()

			select {
			case <-ctx.Done():
				c.Logger().Error("Request timeout")
				return c.String(http.StatusGatewayTimeout, "Request timeout")
			case err := <-done:
				return err
			}
		}
	}
}

// isWebSocketRequest checks if the request is an upgrade request for WebSockets
func isWebSocketRequest(req *http.Request) bool {
	return req.Header.Get("Connection") == "Upgrade" && req.Header.Get("Upgrade") == "websocket"
}
