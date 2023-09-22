package utils

import (
	"github.com/gin-gonic/gin"

	"go-api/pkg/config"
)

// CreateSessionCookie JWT
func CreateSessionCookie(cfg *config.Config, c *gin.Context, session string) {
	c.SetCookie(
		cfg.Session.Name,
		session,
		int(cfg.Session.Duration),
		cfg.Cookie.Path,
		cfg.Cookie.Domain,
		cfg.Cookie.Secure,
		cfg.Cookie.HTTPOnly,
	)
}

// DeleteSessionCookie setting blank
func DeleteSessionCookie(cfg *config.Config, c *gin.Context, sessionName string) {
	c.SetCookie(
		sessionName,
		"",
		-1,
		cfg.Cookie.Path,
		cfg.Cookie.Domain,
		cfg.Cookie.Secure,
		cfg.Cookie.HTTPOnly,
	)
}
