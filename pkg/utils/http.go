package utils

import (
	"github.com/gin-gonic/gin"
)

const HeaderXRequestID = "X-Request-ID"

func GetRequestID(c *gin.Context) string {
	return c.GetHeader(HeaderXRequestID)
}

func GetRemoteAddress(c *gin.Context) string {
	return c.Request.RemoteAddr
}
