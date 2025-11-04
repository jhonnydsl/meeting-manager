package middleware

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Cors() gin.HandlerFunc {
	config := cors.Config{
		AllowOrigins: []string{"http://localhost:4200"},
		AllowMethods: []string{"GET", "Content-Type", "Authorization"},
		AllowHeaders: []string{"Origin", "Content-Type", "Authorization"},
	}
	return cors.New(config)
}