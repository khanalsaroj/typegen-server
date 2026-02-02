package java

import (
	"database/sql"
	"fmt"
	"github.com/khanalsaroj/typegen-server/internal/domain"
	"strings"
)

func scanRows(rows *sql.Rows) ([]domain.SqlData, error) {
	var rowsData []domain.SqlData

	for rows.Next() {
		var dt domain.SqlData
		err := rows.Scan(
			&dt.Ordinal,
			&dt.ColumnName,
			&dt.IsNullable,
			&dt.CharacterMaximumLength,
			&dt.DataType,
			&dt.ColumnKey,
			&dt.ColumnComment,
		)
		if err != nil {
			return nil, err
		}
		rowsData = append(rowsData, dt)
	}

	return rowsData, nil
}

func filterColumns(rowsData []domain.SqlData, skipPrefixes []string,
	excludePrimaryKeys bool) []domain.SqlData {
	var filtered []domain.SqlData

	for _, row := range rowsData {
		if excludePrimaryKeys && strings.Contains(row.ColumnKey, "PRI") {
			continue
		}

		if hasAnyPrefixIgnoreCase(row.ColumnName, skipPrefixes) {
			continue
		}

		filtered = append(filtered, row)
	}

	return filtered
}

func getPrimaryKeys(rowsData []domain.SqlData) []domain.SqlData {
	var primaryKeys []domain.SqlData

	for _, row := range rowsData {
		if strings.Contains(row.ColumnKey, "PRI") {
			primaryKeys = append(primaryKeys, row)
		}
	}

	return primaryKeys
}

func writeColumnList(sb *strings.Builder, columns []domain.SqlData) {
	for i, row := range columns {
		sb.WriteString(fmt.Sprintf("          %s", row.ColumnName))
		if i < len(columns)-1 {
			sb.WriteString(",")
		}
		sb.WriteString("\n")
	}
}

func hasAnyPrefixIgnoreCase(columnName string, prefixes []string) bool {
	columnName = strings.ToLower(columnName)

	for _, prefix := range prefixes {
		if strings.HasPrefix(columnName, strings.ToLower(prefix)) {
			return true
		}
	}
	return false
}

func removeNumberOfLines(sb *strings.Builder, numberOfLines int) {
	xml := sb.String()
	sb.Reset()
	sb.WriteString(xml[:len(xml)-numberOfLines])
}
