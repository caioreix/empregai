package http

import (
	"github.com/gin-gonic/gin"

	"go-api/internal/core/user"
	"go-api/internal/middleware"
)

func MapUserRoutes(group *gin.RouterGroup, h user.Handlers, mw *middleware.Manager) {
	group.POST("/register", h.Register())
	group.POST("/login", h.Login())
	group.POST("/logout", h.Logout())

	group.Use(mw.AuthSession())
	group.GET("/all", h.GetUsers())
	group.GET("/:user_id", h.GetUserByID())
	group.GET("/me", h.GetMe())
	group.PUT("/:user_id", h.Update())
	group.DELETE("/:user_id", h.Delete())
}
