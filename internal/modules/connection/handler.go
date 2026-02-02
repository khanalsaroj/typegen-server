package connection

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/khanalsaroj/typegen-server/internal/domain"
	"github.com/khanalsaroj/typegen-server/internal/pkg/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Create(c *gin.Context) {
	var req DatabaseConnectionsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	connection, err := h.service.Create(c.Request.Context(), &req)
	if err != nil {
		if errors.Is(err, domain.ErrConflict) {
			response.Error(c, http.StatusConflict, "Connection already exists", err)
			return
		}
		response.Error(c, http.StatusInternalServerError, "Failed to create connection", err)
		return
	}

	response.Success(c, http.StatusCreated, "Connection created successfully", connection)
}

func (h *Handler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid connection ID", err)
		return
	}

	connection, err := h.service.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.Error(c, http.StatusNotFound, "Connection not found", err)
			return
		}
		response.Error(c, http.StatusInternalServerError, "Failed to get connection", err)
		return
	}

	response.Success(c, http.StatusOK, "Connection retrieved successfully", connection)
}

func (h *Handler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	users, total, err := h.service.List(c.Request.Context(), page, pageSize)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to list users", err)
		return
	}

	response.SuccessWithPagination(c, http.StatusOK, "Connection retrieved successfully", users, page, pageSize, total)
}

func (h *Handler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid connection ID", err)
		return
	}

	var req DatabaseConnectionsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	connection, err := h.service.Update(c.Request.Context(), uint(id), &req)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.Error(c, http.StatusNotFound, "Connection not found", err)
			return
		}
		response.Error(c, http.StatusInternalServerError, "Failed to update connection", err)
		return
	}

	response.Success(c, http.StatusOK, "Connection updated successfully", connection)
}

func (h *Handler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid connection ID", err)
		return
	}

	if err := h.service.Delete(c.Request.Context(), uint(id)); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.Error(c, http.StatusNotFound, "Connection not found", err)
			return
		}
		response.Error(c, http.StatusInternalServerError, "Failed to delete Connection", err)
		return
	}

	response.Success(c, http.StatusOK, "connection deleted successfully", nil)
}
