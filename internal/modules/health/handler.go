package health

import (
	"github.com/khanalsaroj/typegen-server/internal/domain"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Health(c *gin.Context) {
	result := h.service.GetHealth(c.Request.Context())

	statusCode := http.StatusOK
	if result.Status == domain.Unhealthy {
		statusCode = http.StatusServiceUnavailable
	}
	c.JSON(statusCode, gin.H{
		"status":   result.Status,
		"version":  result.Version,
		"uptime":   result.Uptime,
		"database": result.Database,
	})
}
