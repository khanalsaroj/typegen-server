package mysql

import (
	"database/sql"
	"github.com/khanalsaroj/typegen-server/internal/query"

	"github.com/khanalsaroj/typegen-server/internal/domain"
)

func ReadSchema(info domain.DatabaseConnectionInfo, db *sql.DB, tbN string) (*sql.Rows, error) {
	return db.Query(query.TableColumnDataMySQL, info.DatabaseName, tbN)
}
