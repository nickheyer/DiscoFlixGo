package handlers

import (
	"net/http"

	"github.com/gorilla/sessions"
	echomw "github.com/labstack/echo/v4/middleware"
	"github.com/nickheyer/DiscoFlixGo/config"
	"github.com/nickheyer/DiscoFlixGo/pkg/middleware"
	"github.com/nickheyer/DiscoFlixGo/pkg/services"
)

// BuildRouter builds the router
func BuildRouter(c *services.Container) error {
	// Static files with proper cache control
	// funcmap.File() should be used in templates to append a cache key to the URL in order to break cache
	// after each server restart
	c.Web.Group("", middleware.CacheControl(c.Config.Cache.Expiration.StaticFile)).
		Static(config.StaticPrefix, config.StaticDir)

	// Non-static file route group
	g := c.Web.Group("")

	// Force HTTPS, if enabled
	if c.Config.HTTP.TLS.Enabled {
		g.Use(echomw.HTTPSRedirect())
	}

	g.Use(
		echomw.RemoveTrailingSlashWithConfig(echomw.TrailingSlashConfig{
			RedirectCode: http.StatusMovedPermanently,
		}),
		echomw.Recover(),
		echomw.Secure(),
		echomw.RequestID(),
		middleware.SetLogger(),
		middleware.LogRequest(),
		echomw.Gzip(),
		middleware.Session(sessions.NewCookieStore([]byte(c.Config.App.EncryptionKey))),
		middleware.LoadAuthenticatedUser(c.Auth),
		middleware.ServeCachedPage(c.TemplateRenderer),
		echomw.CSRFWithConfig(echomw.CSRFConfig{
			TokenLookup: "form:csrf",
		}),
	)

	// Error handler
	err := Error{c.TemplateRenderer}
	c.Web.HTTPErrorHandler = err.Page

	// Initialize and register all handlers
	for _, h := range GetHandlers() {
		if err := h.Init(c); err != nil {
			return err
		}

		h.Routes(g)
	}

	return nil
}
