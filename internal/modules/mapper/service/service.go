package service

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/khanalsaroj/typegen-server/internal/domain"
	"github.com/khanalsaroj/typegen-server/internal/modules/connection"
	"github.com/khanalsaroj/typegen-server/internal/modules/helper"
	"github.com/khanalsaroj/typegen-server/internal/modules/mapper/generator"
)

type MprService struct {
	ConnectionService *connection.Service
}

func (s *MprService) Generate(c *gin.Context, req domain.MapperRequest) (string, error) {
	connDetails, err := s.ConnectionService.GetByID(c.Request.Context(), req.ConnectionId)
	if err != nil {
		return "", err
	}
	db, reader, connInfo, err := helper.OpenDatabase(connDetails)
	if err != nil {
		return "", err
	}
	defer func(db *sql.DB) {
		if db.Close() != nil {
		}
	}(db)

	cols, err := reader.ReadSchema(connInfo, db, req.TableName)

	if err != nil {
		return "", err
	}

	mapper, err := generator.NewGenerator(req)
	if err != nil {
		return "", err
	}

	return mapper.Generate(cols, req, req.TableName)
}
