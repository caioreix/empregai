package server

import (
	"net/http"

	"github.com/gin-gonic/gin"

	sessionrepo "go-api/internal/features/session/repository/redisrepo"
	sessionusecase "go-api/internal/features/session/usecase"
	userhandler "go-api/internal/features/user/delivery/http"
	userrepo "go-api/internal/features/user/repository/postgres"
	userusecase "go-api/internal/features/user/usecase"
	"go-api/internal/middleware"
)

func (s *Server) MapHandlers() error {
	// Repository
	userRepo := userrepo.NewUserRepository(s.db)
	sessionRepo := sessionrepo.NewSessionRepository(s.redisClient, s.cfg)

	// UseCase
	userUC := userusecase.NewUserUseCase(s.cfg, userRepo)
	sessionUC := sessionusecase.NewSessionUseCase(sessionRepo, s.cfg)

	// Handler
	userHandlers := userhandler.NewUserHandler(s.cfg, userUC, sessionUC)

	s.gin.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	})

	mw := middleware.New(s.cfg, s.logger, userUC, sessionUC)

	v1 := s.gin.Group("/v1")
	v1.Use(mw.RequestID())

	health := v1.Group("/health")
	authGroup := v1.Group("/auth")

	userhandler.MapUserRoutes(authGroup, userHandlers, mw)

	health.GET("", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	return nil
}
