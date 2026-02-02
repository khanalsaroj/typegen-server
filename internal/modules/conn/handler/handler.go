package http

import (
	"github.com/khanalsaroj/typegen-server/internal/domain"
	"github.com/khanalsaroj/typegen-server/internal/modules/conn/service"
	"github.com/khanalsaroj/typegen-server/internal/modules/connection"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service           service.Service
	connectionService connection.Service
}

func NewHandler(service service.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Connect(c *gin.Context) {
	var req domain.DatabaseConnectionInfo
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, domain.ConnectResponse{
			Message: "Invalid request payload",
			Success: false,
		})
		return
	}

	res := h.service.CheckConnection(req)
	c.JSON(http.StatusOK, res)
}
