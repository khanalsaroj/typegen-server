package response

import (
	"github.com/gin-gonic/gin"
)

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Meta    *Meta       `json:"meta,omitempty"`
}

type Meta struct {
	Page      int   `json:"page"`
	PageSize  int   `json:"page_size"`
	Total     int64 `json:"total"`
	TotalPage int   `json:"total_page"`
}

func Success(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func SuccessWithPagination(c *gin.Context, statusCode int, message string, data interface{}, page, pageSize int, total int64) {
	totalPage := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPage++
	}

	c.JSON(statusCode, Response{
		Success: true,
		Message: message,
		Data:    data,
		Meta: &Meta{
			Page:      page,
			PageSize:  pageSize,
			Total:     total,
			TotalPage: totalPage,
		},
	})
}

func Error(c *gin.Context, statusCode int, message string, err error) {
	resp := Response{
		Success: false,
		Message: message,
	}

	if err != nil {
		resp.Error = err.Error()
	}

	c.JSON(statusCode, resp)
}
