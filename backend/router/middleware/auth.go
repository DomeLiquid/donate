package middleware

// import (
// 	"net/http"
// 	"strings"

// 	"github.com/gin-gonic/gin"
// 	"github.com/lixvyang/dome/pkg/jwt"
// )

// const (
// 	UserIDKey = "uid"
// )

// func JWTNOAuthMiddleware() func(c *gin.Context) {
// 	return func(c *gin.Context) {
// 		authHeader := c.Request.Header.Get("Authorization")
// 		if authHeader == "" {
// 			c.Next()
// 		} else {
// 			parts := strings.SplitN(authHeader, " ", 2)
// 			if !(len(parts) == 2 && parts[0] == "Bearer") {
// 				c.JSON(http.StatusForbidden, gin.H{
// 					"code":    http.StatusForbidden,
// 					"message": "invaild token",
// 				})
// 				c.Abort()
// 				return
// 			}

// 			mc, err := jwt.ParseJwt(parts[1])
// 			if err != nil {
// 				c.JSON(http.StatusForbidden, gin.H{
// 					"code":    http.StatusForbidden,
// 					"message": "invaild token",
// 				})
// 				c.Abort()
// 				return
// 			}
// 			c.Set(UserIDKey, mc.Uid)
// 			c.Next()
// 		}
// 	}
// }

// func JWTAuthMiddleware() func(c *gin.Context) {
// 	return func(c *gin.Context) {
// 		authHeader := c.Request.Header.Get("Authorization")
// 		if authHeader == "" {
// 			c.JSON(http.StatusForbidden, gin.H{
// 				"code":    http.StatusForbidden,
// 				"message": "invaild token",
// 			})
// 			c.Abort()
// 			return
// 		}

// 		parts := strings.SplitN(authHeader, " ", 2)
// 		if !(len(parts) == 2 && parts[0] == "Bearer") {
// 			c.JSON(http.StatusForbidden, gin.H{
// 				"code":    http.StatusForbidden,
// 				"message": "invaild token",
// 			})
// 			c.Abort()
// 			return
// 		}

// 		mc, err := jwt.ParseJwt(parts[1])
// 		if err != nil {
// 			c.JSON(http.StatusForbidden, gin.H{
// 				"code":    http.StatusForbidden,
// 				"message": "invaild token",
// 			})
// 			c.Abort()
// 			return
// 		}

// 		c.Set(UserIDKey, mc.Uid)
// 		c.Next()
// 	}
// }

// /*
// 1. Whether the uid is passed or not, the same API should be used (no restrictions in the middleware).
// 2. Passing uid allows for more optional actions.
// */
// func JWTAuthNotMiddleware() func(c *gin.Context) {
// 	return func(c *gin.Context) {
// 		authHeader := c.Request.Header.Get("Authorization")
// 		if authHeader != "" {
// 			parts := strings.SplitN(authHeader, " ", 2)
// 			if !(len(parts) == 2 && parts[0] == "Bearer") {
// 				c.JSON(http.StatusForbidden, gin.H{
// 					"code":    http.StatusForbidden,
// 					"message": "invaild token",
// 				})
// 				c.Abort()
// 				return
// 			}

// 			mc, err := jwt.ParseJwt(parts[1])
// 			if err != nil {
// 				c.JSON(http.StatusForbidden, gin.H{
// 					"code":    http.StatusForbidden,
// 					"message": "invaild token",
// 				})
// 				c.Abort()
// 				return
// 			}
// 			c.Set(UserIDKey, mc.Uid)
// 		}
// 		c.Next()
// 	}
// }

// /*
// 1. Simulate uid
// */
// func DemoAuthNotMiddleware() func(c *gin.Context) {
// 	return func(c *gin.Context) {
// 		c.Set(UserIDKey, "uid:123123")
// 		c.Next()
// 	}
// }
