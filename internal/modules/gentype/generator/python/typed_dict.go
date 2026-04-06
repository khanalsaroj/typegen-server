package python

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/khanalsaroj/typegen-server/internal/common"
	"github.com/khanalsaroj/typegen-server/internal/domain"
)

type TypedDictDto struct{}

func (d *TypedDictDto) Generate(rows *sql.Rows, req domain.TypeRequest, tbN string, dbType string) (string, error) {

	var sb strings.Builder

	className := req.Prefix + common.ToPascalCase(tbN) + req.Suffix

	var opt domain.PythonTypedDictOptions
	if err := json.Unmarshal(req.Options, &opt); err != nil {
		return "Invalid Python TypedDict Options", fmt.Errorf("invalid typed dict options: %w", err)
	}

	sb.WriteString("from typing import TypedDict")

	if opt.OptionalFields {
		sb.WriteString(", Optional")
	}
	sb.WriteString("\n\n")

	if opt.Total {
		sb.WriteString(fmt.Sprintf("class %s(TypedDict):\n", className))
	} else {
		sb.WriteString(fmt.Sprintf("class %s(TypedDict, total=False):\n", className))
	}

	if opt.Docstrings {
		sb.WriteString(fmt.Sprintf("    \"\"\"%s typed dict\"\"\"\n", className))
	}

	hasField := false

	for rows.Next() {

		var ordinal int
		var columnName string
		var isNullable string
		var characterMaximumLength sql.NullInt16
		var dataType string
		var columnKey string
		var columnComment sql.NullString

		err := rows.Scan(
			&ordinal,
			&columnName,
			&isNullable,
			&characterMaximumLength,
			&dataType,
			&columnKey,
			&columnComment,
		)
		if err != nil {
			return "", err
		}

		hasField = true

		fieldName := common.ToSnakeCase(columnName)
		pyType := mapDBToPythonType(dbType, dataType)

		if opt.OptionalFields || strings.ToLower(isNullable) == "yes" {
			pyType = fmt.Sprintf("Optional[%s]", pyType)
		}

		if opt.Comments && columnComment.Valid && strings.TrimSpace(columnComment.String) != "" {
			sb.WriteString(fmt.Sprintf("    # %s\n", columnComment.String))
		}

		sb.WriteString(fmt.Sprintf("    %s: %s\n", fieldName, pyType))

		if opt.ExtraSpacing {
			sb.WriteString("\n")
		}
	}

	if !hasField {
		sb.WriteString("    pass\n")
	}

	if err := rows.Err(); err != nil {
		return "", err
	}

	return sb.String(), nil
}
