package mysql

import (
	"database/sql"

	"github.com/khanalsaroj/typegen-server/internal/domain"
	"github.com/khanalsaroj/typegen-server/internal/query"
)

func Stats(
	db *sql.DB,
	databaseName string,
	schema string,
) (*domain.Stats, error) {

	stats := &domain.Stats{}

	if err := db.QueryRow(
		query.CountTablesMySQL,
		schema,
	).Scan(&stats.Tables); err != nil {
		return nil, err
	}

	if err := db.QueryRow(
		query.DatabaseSizeMySQL,
		databaseName,
	).Scan(&stats.SizeMB); err != nil {
		return nil, err
	}

	rows, err := db.Query(query.TableNameMySQL, databaseName)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)

	for rows.Next() {
		var t domain.TableInfo
		if err := rows.Scan(&t.Name, &t.ColumnCount); err != nil {
			return nil, err
		}
		stats.TablesInfo = append(stats.TablesInfo, t)
	}

	stats.Tables = len(stats.TablesInfo)
	return stats, nil
}
