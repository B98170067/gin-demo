package middleware

import (
	errno "gin-demo/pkg/error"
	"gin-demo/pkg/response"

	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				response.Error(c, errno.ErrInternal, "internal server error")
				c.Abort()
			}
		}()

		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err

			if appErr, ok := err.(*errno.AppError); ok {
				response.Error(c, appErr.Code, appErr.Message)
				c.Abort()
				return
			}

			response.Error(c, errno.ErrInternal, "internal server error")
			c.Abort()
		}
	}
}
