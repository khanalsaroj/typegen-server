package python

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/khanalsaroj/typegen-server/internal/common"
	"github.com/khanalsaroj/typegen-server/internal/domain"
)

type Dto struct{}

func (d *Dto) Generate(rows *sql.Rows, req domain.TypeRequest, tbN string, dbType string) (string, error) {

	var sb strings.Builder

	className := req.Prefix + common.ToPascalCase(tbN) + req.Suffix

	var opt domain.PythonDataclassOptions
	if err := json.Unmarshal(req.Options, &opt); err != nil {
		return "Invalid Python Options", fmt.Errorf("invalid python options: %w", err)
	}

	// Imports
	sb.WriteString("from dataclasses import dataclass\n")
	if opt.OptionalFields {
		sb.WriteString("from typing import Optional\n")
	}
	sb.WriteString("\n")

	// Dataclass decorator
	decorator := "@dataclass"
	params := []string{}

	if opt.Frozen {
		params = append(params, "frozen=True")
	}
	if opt.Slots {
		params = append(params, "slots=True")
	}
	if opt.KwOnly {
		params = append(params, "kw_only=True")
	}
	if opt.Order {
		params = append(params, "order=True")
	}
	if !opt.Repr {
		params = append(params, "repr=False")
	}
	if !opt.Eq {
		params = append(params, "eq=False")
	}
	if opt.UnsafeHash {
		params = append(params, "unsafe_hash=True")
	}

	if len(params) > 0 {
		decorator = fmt.Sprintf("@dataclass(%s)", strings.Join(params, ", "))
	}

	sb.WriteString(decorator + "\n")

	// Docstring
	if opt.Docstrings {
		sb.WriteString(fmt.Sprintf("class %s:\n", className))
		sb.WriteString(fmt.Sprintf("    \"\"\"%s dataclass\"\"\"\n", className))
	} else {
		sb.WriteString(fmt.Sprintf("class %s:\n", className))
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

		pyType := mapDBToPythonType(dbType, dataType)

		fieldName := common.ToSnakeCase(columnName)

		// Optional handling
		if opt.OptionalFields || strings.ToLower(isNullable) == "yes" {
			pyType = fmt.Sprintf("Optional[%s]", pyType)
		}

		// Comment
		if opt.Comments && columnComment.Valid && strings.TrimSpace(columnComment.String) != "" {
			sb.WriteString(fmt.Sprintf("    # %s\n", columnComment.String))
		}

		// Default values
		defaultVal := ""
		if opt.DefaultValues {
			defaultVal = " = None"
		}

		sb.WriteString(fmt.Sprintf("    %s: %s%s\n", fieldName, pyType, defaultVal))

		if opt.ExtraSpacing {
			sb.WriteString("\n")
		}
	}

	// Handle empty class
	if !hasField {
		sb.WriteString("    pass\n")
	}

	if err := rows.Err(); err != nil {
		return "", err
	}

	return sb.String(), nil
}

func mapDBToPythonType(dbType string, dataType string) string {

	switch strings.ToLower(dbType) {

	case "mysql", "mariadb":
		return mapMySQLToPythonType(dataType)

	case "postgres", "postgresql":
		return mapPostgresToPythonType(dataType)

	case "mssql", "sqlserver":
		return mapMSSQLToPythonType(dataType)

	default:
		return "Any"
	}
}
func mapMySQLToPythonType(t string) string {

	switch strings.ToLower(t) {

	case "int", "integer", "mediumint", "smallint", "tinyint", "bigint":
		return "int"

	case "decimal", "numeric":
		return "float"

	case "float", "double":
		return "float"

	case "varchar", "char", "text", "longtext", "mediumtext", "tinytext":
		return "str"

	case "date":
		return "date"
	case "datetime", "timestamp":
		return "datetime"
	case "time":
		return "time"
	case "year":
		return "int"

	case "boolean", "bit":
		return "bool"

	case "blob", "longblob", "mediumblob", "tinyblob", "binary", "varbinary":
		return "bytes"

	case "json":
		return "dict"

	case "enum":
		return "str"

	default:
		return "Any"
	}
}

func mapPostgresToPythonType(t string) string {

	switch strings.ToLower(t) {

	case "smallint", "int2":
		return "int"

	case "integer", "int", "int4", "serial":
		return "int"

	case "bigint", "int8", "bigserial":
		return "int"

	case "decimal", "numeric", "money":
		return "float"

	case "real", "float4", "double precision", "float8":
		return "float"

	case "varchar", "character varying", "char", "character", "text", "citext":
		return "str"

	case "date":
		return "date"
	case "time", "time without time zone":
		return "time"
	case "timestamp", "timestamp without time zone":
		return "datetime"
	case "timestamptz", "timestamp with time zone":
		return "datetime"

	case "boolean", "bool":
		return "bool"

	case "bytea":
		return "bytes"

	case "uuid":
		return "str"

	case "json", "jsonb":
		return "dict"

	// arrays (basic)
	case "_int4", "integer[]":
		return "list[int]"
	case "_int8", "bigint[]":
		return "list[int]"
	case "_text", "text[]":
		return "list[str]"

	default:
		return "Any"
	}
}

func mapMSSQLToPythonType(t string) string {

	switch strings.ToLower(t) {

	case "int", "bigint", "smallint", "tinyint":
		return "int"

	case "decimal", "numeric", "money", "smallmoney":
		return "float"

	case "float", "real":
		return "float"

	case "varchar", "nvarchar", "char", "nchar", "text", "ntext":
		return "str"

	case "date":
		return "date"
	case "time":
		return "time"
	case "datetime", "datetime2", "smalldatetime":
		return "datetime"
	case "datetimeoffset":
		return "datetime"

	case "bit":
		return "bool"

	case "binary", "varbinary", "image", "rowversion", "timestamp":
		return "bytes"

	case "uniqueidentifier":
		return "str"

	default:
		return "Any"
	}
}
