package middleware

import (
	"context"

	"github.com/gin-gonic/gin"

	"go-api/internal/core/user"
	"go-api/pkg/apierrors"
	"go-api/pkg/logger"
	"go-api/pkg/utils"
)

func (m *Manager) AuthSession() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := utils.GetRequestID(c)

		sessionID, err := c.Cookie(m.cfg.Session.Name)
		if err != nil {
			m.log.Warn("Failed getting session cookie in auth session middleware", logger.Fields{
				"err":        err,
				"request_id": requestID,
			})

			c.JSON(apierrors.Unauthorized("", "").JSON())
			c.Abort()
			return
		}

		session, err := m.sessionUC.GetSessionByID(c.Request.Context(), sessionID)
		if err != nil {
			m.log.Warn("Failed getting session by id in auth session middleware", logger.Fields{
				"err":        err,
				"request_id": requestID,
			})

			c.JSON(apierrors.Unauthorized("", "").JSON())
			c.Abort()
			return
		}

		usr, err := m.userUC.GetByID(c.Request.Context(), session.UserID)
		if err != nil {
			m.log.Warn("Failed getting user by id in auth session middleware", logger.Fields{
				"err":        err,
				"request_id": requestID,
			})

			c.JSON(apierrors.Unauthorized("", "").JSON())
			c.Abort()
			return
		}

		ctx := context.WithValue(c.Request.Context(), user.CtxKey{}, usr)
		c.Request = c.Request.WithContext(ctx)

		m.log.Info("Succeeded auth session middleware", logger.Fields{
			"request_id":     requestID,
			"remote_address": utils.GetRemoteAddress(c),
			"user_id":        usr.ID,
			"session_id":     sessionID,
		})

		c.Next()
	}
}
