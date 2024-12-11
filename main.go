package main

import (
	"booking-event-server/config"
	"booking-event-server/router"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode)

	config.LoadDB()

	r := gin.New()
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	api := r.Group("/api")
	r.SetTrustedProxies([]string{"127.0.0.1"})

	router.AuthRouter(api)

	port := os.Getenv("PORT")
	if port == "" {
		port = "5001"
	}

	fmt.Printf("Server running on: http://localhost:%v\n", port)
	r.Run(fmt.Sprintf(":%v", port))
}
