package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jhonnydsl/gerenciamento-de-reunioes/src/dtos"
)

func ErrorMiddlewareHandle() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err

			apiErr, ok := err.(*dtos.APIError)
			if ok {
				c.JSON(apiErr.StatusCode, gin.H{"error": apiErr.Message})
				c.Abort()
				return 
			}

			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error internal Server"})
			c.Abort()
		}
	}
}