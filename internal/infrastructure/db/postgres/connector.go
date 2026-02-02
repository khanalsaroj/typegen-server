package postgres

import (
	"database/sql"
	"fmt"
	"github.com/khanalsaroj/typegen-server/internal/domain"

	_ "github.com/lib/pq"
)

type Connector struct{}

func (c *Connector) Open(req domain.DatabaseConnectionInfo) (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		req.Host,
		req.Port,
		req.Username,
		req.Password,
		req.DatabaseName,
	)

	return sql.Open("postgres", dsn)
}

func (c *Connector) Stats(db *sql.DB, databaseName string, schema string) (*domain.Stats, error) {
	return Stats(db, databaseName, schema)
}

func (c *Connector) ReadSchema(req domain.DatabaseConnectionInfo, db *sql.DB, tbN string) (*sql.Rows, error) {
	return ReadSchema(req, db, tbN)
}
