package connector

import (
	"database/sql"
	"github.com/khanalsaroj/typegen-server/internal/domain"
)

type DBConnector interface {
	Open(req domain.DatabaseConnectionInfo) (*sql.DB, error)
	Stats(db *sql.DB, databaseName, schema string) (*domain.Stats, error)
	ReadSchema(req domain.DatabaseConnectionInfo, db *sql.DB, tbN string) (*sql.Rows, error)
}
