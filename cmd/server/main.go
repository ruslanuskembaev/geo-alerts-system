package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/ruslanuskembaev/geo-alerts-system/internal/config"
)

func main() {
	cfg := config.Load()

	fmt.Printf("Starting server on port %s\n", cfg.ServerPort)

	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.Run(":" + cfg.ServerPort)
}
