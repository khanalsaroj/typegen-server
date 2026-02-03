package typescript

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/khanalsaroj/typegen-server/internal/common"
	"github.com/khanalsaroj/typegen-server/internal/domain"
)

type Dto struct{}

func (r *Dto) Generate(rows *sql.Rows, info domain.TypeRequest, tbN string, dbType string) (string, error) {
	var sb strings.Builder

	style := info.Style
	tableName := info.Prefix + common.ToPascalCase(tbN) + info.Suffix

	var opt domain.TypeScriptOptions
	if err := json.Unmarshal(info.Options, &opt); err != nil {
		return "Invalid Java Options", fmt.Errorf("invalid TypeScript options: %w", err)
	}

	export := ""
	if opt.ExportAllTypes {
		if style == "class" {
			export = "export default"
		} else {
			export = "export"
		}
	}

	switch style {
	case "interface":
		sb.WriteString(fmt.Sprintf("%s interface %s {\n", export, tableName))
	case "class":
		sb.WriteString(fmt.Sprintf("%s class %s {\n", export, tableName))
	default:
		sb.WriteString(fmt.Sprintf(" %s type %s = {\n", export, tableName))
	}

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

		var tsType string

		switch strings.ToLower(dbType) {
		case "mysql":
			tsType = mapMySQLToTSType(dataType)
		case "postgres":
			tsType = mapPostgresqlToTSType(dataType)
		case "mssql", "sqlserver", "sql_server":
			tsType = mapMSSQLToTSType(dataType)
		default:
			tsType = "any"
		}

		optional := ""
		if opt.OptionalProperties {
			optional = "?"
		} else {
			if isNullable == "YES" {
				optional = "?"
			}
		}

		fieldName := common.ToCamelCase(columnName)

		readonly := ""
		if opt.ReadonlyProperties {
			readonly = "readonly "
		}

		if opt.JSDocComments {
			sb.WriteString(fmt.Sprintf("  /**@type {%s} */\n", tsType))
		}

		if opt.Comments {
			if columnComment.Valid && strings.TrimSpace(columnComment.String) != "" {
				sb.WriteString(fmt.Sprintf(
					"  /** %s */\n",
					columnComment.String,
				))
			}
		}

		sb.WriteString(fmt.Sprintf(
			"  %s%s%s: %s\n",
			readonly,
			fieldName,
			optional,
			tsType,
		))
		if opt.ExtraSpacing {
			sb.WriteString(fmt.Sprintf("\n"))
		}

	}

	sb.WriteString("}\n")

	if err := rows.Err(); err != nil {
		return "", err
	}

	return sb.String(), nil
}

func mapPostgresqlToTSType(udtName string) string {
	isArray := false
	baseType := udtName

	if strings.HasPrefix(udtName, "_") {
		isArray = true
		baseType = strings.TrimPrefix(udtName, "_")
	}
	var tsType string
	switch baseType {
	case "int2", "int4", "int8", "oid":
		tsType = "number"
	case "float4", "float8", "numeric", "money":
		tsType = "number"
	case "bool":
		tsType = "boolean"
	case "text", "varchar", "bpchar", "char", "name",
		"uuid", "inet", "cidr", "macaddr":
		tsType = "string"
	case "date", "timestamp", "timestamptz", "time", "timetz":
		tsType = "string | Date"
	case "json", "jsonb":
		tsType = "Record<string, any>"
	case "bytea":
		tsType = "Uint8Array"
	case "int4range", "int8range", "numrange",
		"tsrange", "tstzrange", "daterange":
		tsType = "string"
	default:
		tsType = "any"
	}
	if isArray {
		tsType = tsType + "[]"
	}

	return tsType
}

func mapMySQLToTSType(mysqlType string) string {
	switch strings.ToLower(mysqlType) {

	case "int", "bigint", "smallint", "tinyint", "decimal", "float", "double", "long":
		return "number"

	case "varchar", "char", "text", "longtext":
		return "string"

	case "datetime", "timestamp", "date":
		return "string | Date"

	case "boolean", "bit":
		return "boolean"

	default:
		return "any"
	}
}

func mapMSSQLToTSType(mssqlType string) string {
	switch strings.ToLower(mssqlType) {
	case "int", "bigint", "smallint", "tinyint",
		"decimal", "numeric", "float", "real",
		"money", "smallmoney":
		return "number"
	case "varchar", "nvarchar", "char", "nchar", "text", "ntext":
		return "string"
	case "datetime", "datetime2", "smalldatetime", "date", "time", "datetimeoffset":
		return "string | Date"
	case "bit":
		return "boolean"
	case "binary", "varbinary", "image":
		return "Uint8Array"
	case "uniqueidentifier":
		return "string"
	case "xml":
		return "string"
	case "sql_variant", "hierarchyid", "geometry", "geography":
		return "any"
	case "rowversion", "timestamp":
		return "Uint8Array"
	default:
		return "any"
	}
}
