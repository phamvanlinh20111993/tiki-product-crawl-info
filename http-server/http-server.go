package http_server

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func httpServer() {
	// Create a Gin router with default middleware (logger and recovery)
	router := gin.Default()

	// Define a simple GET endpoint
	router.GET("/ping", func(c *gin.Context) {
		// Return JSON response
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	// Start server on port 8080 (default)
	// Server will listen on 0.0.0.0:8080 (localhost:8080 on Windows)
	if err := router.Run(); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}

	// Listen and serve on 0.0.0.0:8080
	err := router.Run(":8080")
	if err != nil {
		return
	}
}

var HttpServer = httpServer
