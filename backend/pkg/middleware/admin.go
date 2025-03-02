package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
)

type Admin struct {
	ak string
	sk string
}

var _admin = &Admin{}

func InitAdmin(ak, sk string) {
	_admin = &Admin{
		ak: ak,
		sk: sk,
	}
}

func (a *Admin) Allow(ak, sk string) bool {
	return a.ak == ak && a.sk == sk
}

func AdminAuthMiddleware(enable bool) func(c *gin.Context) {
	return func(c *gin.Context) {
		if !enable {
			c.Next()
			return
		}

		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			c.JSON(403, gin.H{
				"code":    403,
				"message": "invaild admin token",
			})
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(403, gin.H{
				"code":    403,
				"message": "invaild admin token",
			})
			c.Abort()
			return
		}
		split := strings.SplitN(parts[1], ":", 2)
		if !(len(split) == 2) {
			c.JSON(403, gin.H{
				"code":    403,
				"message": "invaild admin token",
			})
			c.Abort()
			return
		}

		if !_admin.Allow(split[0], split[1]) {
			c.JSON(403, gin.H{
				"code":    403,
				"message": "invaild admin token",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
