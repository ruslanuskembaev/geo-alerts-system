package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

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
		log.Println("No .env file found, using defaults")
	} else {
		log.Println(".env file loaded")
	}

	cfg := config.Load()
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	fmt.Printf("\nStarting Geo Alerts System\n")
	fmt.Printf("   Server Port: %s\n", cfg.ServerPort)
	fmt.Printf("   API Key: %s...\n", cfg.APIKey[:min(10, len(cfg.APIKey))])
	fmt.Println()

	// Инициализация инфраструктуры
	dbPool, err := repository.NewPostgresPool(cfg)
	if err != nil {
		log.Fatal("Failed to connect to PostgreSQL:", err)
	}

	redisClient := repository.NewRedisClient(cfg)

	// Инициализация слоёв
	incidentRepo := repository.NewIncidentRepository(dbPool)
	checkRepo := repository.NewLocationCheckRepository(dbPool)
	cache := repository.NewIncidentCache(redisClient, cfg.CacheTTL)
	queue := repository.NewWebhookQueue(redisClient)
	systemRepo := repository.NewSystemRepository(dbPool, redisClient)

	incidentService := service.NewIncidentService(incidentRepo, cache, checkRepo)
	locationService := service.NewLocationService(incidentRepo, cache, checkRepo, queue)
	healthService := service.NewHealthService(systemRepo, cfg.HealthTimeout)

	webhookSender := service.NewWebhookSender(cfg.WebhookURL, cfg.WebhookTimeout)
	webhookWorker := service.NewWebhookWorker(queue, webhookSender, cfg.WebhookRetryAttempts, cfg.WebhookRetryDelay)
	workerCtx, cancelWorker := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		webhookWorker.Start(workerCtx)
	}()

	incidentHandler := handler.NewIncidentHandler(incidentService, cfg.StatsTimeWindow)
	locationHandler := handler.NewLocationHandler(locationService)
	healthHandler := handler.NewHealthHandler(healthService)

	// HTTP сервер
	r := gin.Default()

	// API v1
	api := r.Group("/api/v1")
	{
		// Health check (публичный)
		api.GET("/system/health", healthHandler.Health)

		// Location check (публичный)
		api.POST("/location/check", locationHandler.Check)

		// Incidents (защищённые endpoints)
		incidents := api.Group("/incidents")
		incidents.Use(handler.AuthMiddleware(cfg.APIKey))
		{
			incidents.POST("", incidentHandler.Create)
			incidents.GET("", incidentHandler.List)
			incidents.GET("/stats", incidentHandler.Stats)
			incidents.GET("/:id", incidentHandler.GetByID)
			incidents.PUT("/:id", incidentHandler.Update)
			incidents.DELETE("/:id", incidentHandler.Delete)
		}
	}

	fmt.Println("Available endpoints:")
	fmt.Println("   GET  /api/v1/system/health          (public)")
	fmt.Println("   POST /api/v1/location/check         (public)")
	fmt.Println("   POST /api/v1/incidents              (protected)")
	fmt.Println("   GET  /api/v1/incidents              (protected)")
	fmt.Println("   GET  /api/v1/incidents/stats         (protected)")
	fmt.Println("   GET  /api/v1/incidents/:id          (protected)")
	fmt.Println("   PUT  /api/v1/incidents/:id          (protected)")
	fmt.Println("   DELETE /api/v1/incidents/:id        (protected)")
	fmt.Println()
	fmt.Printf("Server running at http://localhost:%s\n\n", cfg.ServerPort)

	server := &http.Server{
		Addr:              ":" + cfg.ServerPort,
		Handler:           r,
		ReadTimeout:       cfg.HTTPReadTimeout,
		ReadHeaderTimeout: cfg.HTTPReadTimeout,
		WriteTimeout:      cfg.HTTPWriteTimeout,
		IdleTimeout:       cfg.HTTPIdleTimeout,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Failed to start server:", err)
		}
	}()

	<-ctx.Done()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("HTTP server shutdown error: %v\n", err)
	}
	cancelWorker()
	wg.Wait()
	_ = redisClient.Close()
	dbPool.Close()
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
