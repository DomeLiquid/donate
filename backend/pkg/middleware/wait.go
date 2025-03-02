package middleware

import (
	"github.com/gin-gonic/gin"
)

func WaitMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		// TODO Wait
		// time.Sleep(2 * time.Second)
		c.Next()
	}
}
