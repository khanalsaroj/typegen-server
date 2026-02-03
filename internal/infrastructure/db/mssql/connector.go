package mssql

import (
	"database/sql"
	"fmt"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/khanalsaroj/typegen-server/internal/domain"
)

type Connector struct{}

func (c *Connector) Open(req domain.DatabaseConnectionInfo) (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"sqlserver://%s:%s@%s:%d?database=%s&encrypt=disable",
		req.Username,
		req.Password,
		req.Host,
		req.Port,
		req.DatabaseName,
	)
	return sql.Open("sqlserver", dsn)
}

func (c *Connector) Stats(db *sql.DB, _, schema string) (*domain.Stats, error) {
	return Stats(db, schema)
}

func (c *Connector) ReadSchema(req domain.DatabaseConnectionInfo, db *sql.DB, tbN string) (*sql.Rows, error) {
	return ReadSchema(req, db, tbN)
}
