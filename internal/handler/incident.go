package handler

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/ruslanuskembaev/geo-alerts-system/internal/domain"
	"github.com/ruslanuskembaev/geo-alerts-system/internal/repository"
	"github.com/ruslanuskembaev/geo-alerts-system/internal/service"
)

// IncidentHandler обработчик HTTP запросов для инцидентов
type IncidentHandler struct {
	service     *service.IncidentService
	statsWindow time.Duration
}

func NewIncidentHandler(service *service.IncidentService, statsWindow time.Duration) *IncidentHandler {
	return &IncidentHandler{
		service:     service,
		statsWindow: statsWindow,
	}
}

// Create создаёт новый инцидент
func (h *IncidentHandler) Create(c *gin.Context) {
	var req domain.CreateIncidentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid request",
			"details": err.Error(),
		})
		return
	}

	incident, err := h.service.Create(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, incident)
}

// GetByID получает инцидент по ID
func (h *IncidentHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	incident, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, repository.ErrNotFound) {
			status = http.StatusNotFound
			c.JSON(status, gin.H{"error": "incident not found"})
			return
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, incident)
}

// List возвращает список всех активных инцидентов
func (h *IncidentHandler) List(c *gin.Context) {
	page, pageSize, limit, offset := parsePagination(c)

	incidents, total, err := h.service.List(c.Request.Context(), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"incidents": incidents,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// Update обновляет инцидент
func (h *IncidentHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var req domain.UpdateIncidentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid request",
			"details": err.Error(),
		})
		return
	}

	incident, err := h.service.Update(c.Request.Context(), id, req)
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, repository.ErrNotFound) {
			status = http.StatusNotFound
			c.JSON(status, gin.H{"error": "incident not found"})
			return
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, incident)
}

// Delete деактивирует инцидент
func (h *IncidentHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	if err := h.service.Deactivate(c.Request.Context(), id); err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, repository.ErrNotFound) {
			status = http.StatusNotFound
			c.JSON(status, gin.H{"error": "incident not found"})
			return
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "incident deactivated successfully",
	})
}

// Stats возвращает статистику по инцидентам
func (h *IncidentHandler) Stats(c *gin.Context) {
	since := time.Now().Add(-h.statsWindow)

	stats, err := h.service.StatsByIncident(c.Request.Context(), since)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"window_minutes": int(h.statsWindow.Minutes()),
		"stats":          stats,
	})
}

func parsePagination(c *gin.Context) (int, int, int, int) {
	page := 1
	pageSize := 20
	const maxPageSize = 100

	if value := c.Query("page"); value != "" {
		if parsed, err := strconv.Atoi(value); err == nil && parsed > 0 {
			page = parsed
		}
	}
	if value := c.Query("page_size"); value != "" {
		if parsed, err := strconv.Atoi(value); err == nil && parsed > 0 {
			pageSize = parsed
		}
	}
	if pageSize > maxPageSize {
		pageSize = maxPageSize
	}

	limit := pageSize
	offset := (page - 1) * pageSize
	return page, pageSize, limit, offset
}
