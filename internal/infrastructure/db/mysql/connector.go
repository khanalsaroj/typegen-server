package mysql

import (
	"database/sql"
	"fmt"
	"github.com/khanalsaroj/typegen-server/internal/domain"

	_ "github.com/go-sql-driver/mysql"
)

type Connector struct{}

func (c *Connector) Open(req domain.DatabaseConnectionInfo) (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=UTC",
		req.Username,
		req.Password,
		req.Host,
		req.Port,
		req.DatabaseName,
	)
	return sql.Open("mysql", dsn)
}

func (c *Connector) Stats(db *sql.DB, databaseName, schema string) (*domain.Stats, error) {
	return Stats(db, databaseName, schema)
}

func (c *Connector) ReadSchema(req domain.DatabaseConnectionInfo, db *sql.DB, tbN string) (*sql.Rows, error) {
	return ReadSchema(req, db, tbN)
}
