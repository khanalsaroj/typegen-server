package mssql

import (
	"database/sql"

	"github.com/khanalsaroj/typegen-server/internal/domain"
	"github.com/khanalsaroj/typegen-server/internal/query"
)

func Stats(
	db *sql.DB,
	schema string,
) (*domain.Stats, error) {

	stats := &domain.Stats{}

	if err := db.QueryRow(
		query.CountTablesMSSQL,
		schema,
	).Scan(&stats.Tables); err != nil {
		return nil, err
	}
	if err := db.QueryRow(
		query.DatabaseSizeMSSQL,
		schema,
	).Scan(&stats.SizeMB); err != nil {
		return nil, err
	}
	rows, err := db.Query(query.TableNameMSSQL, schema)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		if rows.Close() != nil {
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
