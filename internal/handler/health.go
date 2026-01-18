package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/ruslanuskembaev/geo-alerts-system/internal/service"
)

// HealthHandler обработчик для health check
type HealthHandler struct {
	service *service.HealthService
}

func NewHealthHandler(service *service.HealthService) *HealthHandler {
	return &HealthHandler{service: service}
}

// Health возвращает статус системы
func (h *HealthHandler) Health(c *gin.Context) {
	health := h.service.Check(c.Request.Context())
	status := http.StatusOK
	if health.Status != "healthy" {
		status = http.StatusServiceUnavailable
	}
	c.JSON(status, health)
}
