package connector

import (
	"errors"
	"github.com/khanalsaroj/typegen-server/internal/domain"
	"github.com/khanalsaroj/typegen-server/internal/infrastructure/db/mysql"
	"github.com/khanalsaroj/typegen-server/internal/infrastructure/db/postgres"
)

func New(req domain.DatabaseConnectionInfo) (DBConnector, error) {
	switch req.DbType {
	case "mysql":
		return &mysql.Connector{}, nil
	case "postgres":
		return &postgres.Connector{}, nil
	default:
		return nil, errors.New("unsupported database")
	}
}
