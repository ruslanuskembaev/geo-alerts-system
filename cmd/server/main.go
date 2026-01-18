package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/ruslanuskembaev/geo-alerts-system/internal/config"
	"github.com/ruslanuskembaev/geo-alerts-system/internal/handler"
	"github.com/ruslanuskembaev/geo-alerts-system/internal/repository"
	"github.com/ruslanuskembaev/geo-alerts-system/internal/service"
)

func main() {
	// Загрузка .env
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  No .env file found, using defaults")
	} else {
		log.Println("✅ .env file loaded")
	}

	cfg := config.Load()

	fmt.Printf("\n🚀 Starting Geo Alerts System\n")
	fmt.Printf("   Server Port: %s\n", cfg.ServerPort)
	fmt.Printf("   API Key: %s...\n", cfg.APIKey[:min(10, len(cfg.APIKey))])
	fmt.Println()

	// Инициализация слоёв
	incidentRepo := repository.NewIncidentRepository()
	incidentService := service.NewIncidentService(incidentRepo)
	incidentHandler := handler.NewIncidentHandler(incidentService)
	healthHandler := handler.NewHealthHandler()

	// HTTP сервер
	r := gin.Default()

	// API v1
	api := r.Group("/api/v1")
	{
		// Health check (публичный)
		api.GET("/system/health", healthHandler.Health)

		// Incidents (защищённые endpoints)
		incidents := api.Group("/incidents")
		incidents.Use(handler.AuthMiddleware(cfg.APIKey))
		{
			incidents.POST("", incidentHandler.Create)
			incidents.GET("", incidentHandler.List)
			incidents.GET("/:id", incidentHandler.GetByID)
			incidents.PUT("/:id", incidentHandler.Update)
			incidents.DELETE("/:id", incidentHandler.Delete)
		}
	}

	fmt.Println("📋 Available endpoints:")
	fmt.Println("   GET  /api/v1/system/health          (public)")
	fmt.Println("   POST /api/v1/incidents              (protected)")
	fmt.Println("   GET  /api/v1/incidents              (protected)")
	fmt.Println("   GET  /api/v1/incidents/:id          (protected)")
	fmt.Println("   PUT  /api/v1/incidents/:id          (protected)")
	fmt.Println("   DELETE /api/v1/incidents/:id        (protected)")
	fmt.Println()
	fmt.Printf("✅ Server running at http://localhost:%s\n\n", cfg.ServerPort)

	if err := r.Run(":" + cfg.ServerPort); err != nil {
		log.Fatal("❌ Failed to start server:", err)
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
