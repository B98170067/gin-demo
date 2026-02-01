package middleware

import (
	errno "gin-demo/pkg/error"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("secret")

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := strings.Replace(c.GetHeader("Authorization"), "Bearer ", "", 1)
		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(200, gin.H{
				"code":    errno.ErrUnauthorized,
				"message": "unauthorized",
			})
			return
		}
		c.Next()
	}
}
