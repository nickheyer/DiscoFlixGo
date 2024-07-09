package services

import (
	"net/http/httptest"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/websocket"
)

// helper function to start Echo server for testing
func startTestServer(t *testing.T, wsClient *WsClient) *httptest.Server {
	e := wsClient.GetEchoInstance()
	ts := httptest.NewServer(e)
	t.Cleanup(ts.Close)
	return ts
}

// helper function to connect to WebSocket
func connectWebSocket(t *testing.T, url string) *websocket.Conn {
	ws, err := websocket.Dial(url, "ws", "http://localhost")
	require.NoError(t, err)
	return ws
}

func TestWsClient_AddAndRemoveConnections(t *testing.T) {
	e := echo.New()
	wsClient := NewWsClient(e)
	e.GET("/ws", wsClient.AddConnectionToPool)
	ts := startTestServer(t, wsClient)
	wsURL := "ws" + ts.URL[len("http"):]

	// Connect and disconnect a WebSocket client
	ws := connectWebSocket(t, wsURL+"/ws")
	defer ws.Close()

	// Allow some time for the connection to be registered
	time.Sleep(100 * time.Millisecond)
	assert.Equal(t, 1, len(wsClient.GetConnections()))

	// Close the WebSocket and allow some time for removal
	ws.Close()
	wsClient.Remove(ws)
	time.Sleep(100 * time.Millisecond)
	assert.Equal(t, 0, len(wsClient.GetConnections()))
}
