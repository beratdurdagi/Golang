package main

import (
	routes "JWT_PROJECT/routes"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8000"
	}

	router := gin.New()

	router.Use(gin.Logger())

	routes.AuthRoutes(router)
	routes.UserRoutes(router)

	router.GET("/api-v1", func(c *gin.Context) {
		c.JSON(http.StatusInternalServerError, gin.H{"success": "true", "message": "Access granted for apiv1"})

	})

	router.GET("/api-v2", func(c *gin.Context) {
		c.JSON(http.StatusInternalServerError, gin.H{"success": "true", "message": "Access granted for v2"})
	})
	router.Run(":" + port)
}
