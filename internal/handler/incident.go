package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/ruslanuskembaev/geo-alerts-system/internal/domain"
	"github.com/ruslanuskembaev/geo-alerts-system/internal/service"
)

// IncidentHandler обработчик HTTP запросов для инцидентов
type IncidentHandler struct {
	service *service.IncidentService
}

func NewIncidentHandler(service *service.IncidentService) *IncidentHandler {
	return &IncidentHandler{service: service}
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

	incident, err := h.service.Create(req)
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

	incident, err := h.service.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "incident not found",
		})
		return
	}

	c.JSON(http.StatusOK, incident)
}

// List возвращает список всех активных инцидентов
func (h *IncidentHandler) List(c *gin.Context) {
	incidents, err := h.service.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"incidents": incidents,
		"total":     len(incidents),
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

	incident, err := h.service.Update(id, req)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "incident not found",
		})
		return
	}

	c.JSON(http.StatusOK, incident)
}

// Delete деактивирует инцидент
func (h *IncidentHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	if err := h.service.Delete(id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "incident not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "incident deactivated successfully",
	})
}
