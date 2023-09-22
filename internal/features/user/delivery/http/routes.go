package http

import (
	"github.com/gin-gonic/gin"

	"go-api/internal/core/user"
)

func MapUserRoutes(group *gin.RouterGroup, h user.Handlers) {
	group.POST("/register", h.Register())
	group.POST("/login", h.Login())
	group.POST("/logout", h.Logout())
	group.GET("/all", h.GetUsers())
	group.GET("/:user_id", h.GetUserByID())
	// TODO Implement middleware to next routes
	// group.Use(middleware.AuthJWT)
	// group.Use(middleware.AuthSession)

	group.GET("/me", h.GetMe())
	group.PUT("/:user_id", h.Update())
	group.DELETE("/:user_id", h.Delete())
}
