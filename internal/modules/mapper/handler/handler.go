package handler

import (
	"net/http"

	"github.com/khanalsaroj/typegen-server/internal/domain"
	"github.com/khanalsaroj/typegen-server/internal/modules/mapper/service"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *service.MprService
}

func New(service *service.MprService) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) GenerateMapper(c *gin.Context) {
	var req domain.MapperRequest

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
