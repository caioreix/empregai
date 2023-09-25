package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"go-api/pkg/utils"
)

// New initializes the RequestID middleware.
func (m *Manager) RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		rid := c.GetHeader(utils.HeaderXRequestID)
		if rid == "" {
			rid = uuid.New().String()
			c.Request.Header.Add(utils.HeaderXRequestID, rid)
		}

		c.Header(utils.HeaderXRequestID, rid)
		c.Next()
	}
}
