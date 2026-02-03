package mssql

import (
	"database/sql"

	"github.com/khanalsaroj/typegen-server/internal/domain"
	"github.com/khanalsaroj/typegen-server/internal/query"
)

func ReadSchema(info domain.DatabaseConnectionInfo, db *sql.DB, tbN string) (*sql.Rows, error) {
	return db.Query(query.TableColumnDataMSSQL, info.SchemaName, tbN, info.DatabaseName)
}
