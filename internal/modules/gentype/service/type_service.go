package service

import (
	"database/sql"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/khanalsaroj/typegen-server/internal/domain"
	"github.com/khanalsaroj/typegen-server/internal/modules/connection"
	"github.com/khanalsaroj/typegen-server/internal/modules/gentype/generator"
	"github.com/khanalsaroj/typegen-server/internal/modules/helper"
)

type TypeService struct {
	ConnectionService *connection.Service
}

func (s *TypeService) Generate(c *gin.Context, req domain.TypeRequest) (string, error) {

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

	var result strings.Builder
	for _, value := range req.TableNames {
		cols, err := reader.ReadSchema(connInfo, db, value)
		if err != nil {
			return "", err
		}

		generator, err := gen.NewGenerator(req)
		if err != nil {
			return "", err
		}

		output, err := generator.Generate(cols, req, value, connInfo.DbType)
		if err != nil {
			return "", err
		}
		result.WriteString(output)
		result.WriteString("\n")
	}
	return result.String(), nil
}
