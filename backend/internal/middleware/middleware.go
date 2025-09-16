package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

// CORSMiddleware handles CORS headers
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// LoggerMiddleware provides custom logging
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		
		// Process request
		c.Next()
		
		// Log after request is processed
		end := time.Now()
		latency := end.Sub(start)
		
		// Log the request details
		statusCode := c.Writer.Status()
		logMessage := fmt.Sprintf("[%s] %s %s %s %s %d %s %s\n",
			end.Format("2006-01-02 15:04:05"),
			c.ClientIP(),
			c.Request.Method,
			c.Request.URL.Path,
			c.Request.Proto,
			statusCode,
			latency.String(),
			c.Request.UserAgent(),
		)
		
		gin.DefaultWriter.Write([]byte(logMessage))
	}
}