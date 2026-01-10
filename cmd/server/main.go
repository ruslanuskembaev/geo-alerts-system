package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/ruslanuskembaev/geo-alerts-system/internal/config"
)

func main() {

	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found, proceeding with environment variables")
	}

	cfg := config.Load()
	fmt.Printf("Starting server on port %s\n", cfg.ServerPort)
	r := gin.Default()

	r.GET("/secure", func(c *gin.Context) {
		apiKey := c.GetHeader("X-API-KEY")
		if apiKey != cfg.APIKey {
			c.JSON(401, gin.H{"error": "Unauthorized"})
			return
		}
		c.JSON(200, gin.H{"message": "Authorized access"})
	})

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.Run(":" + cfg.ServerPort)
}
