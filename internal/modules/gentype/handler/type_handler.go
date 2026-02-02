package handler

import (
	"net/http"

	"github.com/khanalsaroj/typegen-server/internal/domain"
	"github.com/khanalsaroj/typegen-server/internal/modules/gentype/service"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *service.TypeService
}

func New(service *service.TypeService) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) GenerateType(c *gin.Context) {
	var req domain.TypeRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	result, err := h.service.Generate(c, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.String(http.StatusOK, result)
}
