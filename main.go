package main

import (
	"booking-event-server/config"
	"booking-event-server/docs"
	"booking-event-server/middleware"
	"booking-event-server/router"
	"fmt"
	"os"

	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/gin-gonic/gin"
)

// @BasePath /api/v1

// PingExample godoc
// @Summary ping example
// @Schemes
// @Description do ping
// @Tags example
// @Accept json
// @Produce json
// @Success 200 {string} Helloworld
// @Router /example/helloworld [get]

func main() {
	gin.SetMode(gin.ReleaseMode)

	config.LoadDB()

	r := gin.New()
	r.Use(middleware.CorsSetting())

	docs.SwaggerInfo.BasePath = "/"
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"success": true,
			"message": "welcome to RPB Api",
			"version": "1.0.0",
		})
	})

	url := ginSwagger.URL("http://localhost:5000/swagger/doc.json")
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler, url))

	api := r.Group("/api")

	router.AuthRouter(api)
	router.EventRouter(api)

	port := os.Getenv("PORT")
	if port == "" {
		port = "5001"
	}

	fmt.Printf("Server running on: http://localhost:%v\n", port)
	r.Run(fmt.Sprintf(":%v", port))
}
