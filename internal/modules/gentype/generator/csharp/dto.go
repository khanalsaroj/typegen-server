package csharp

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

	tableName := req.Prefix + common.ToPascalCase(tbN) + req.Suffix

	var opt domain.CSharpDtoOptions
	if err := json.Unmarshal(req.Options, &opt); err != nil {
		return "Invalid CSharp Options", fmt.Errorf("invalid CSharp options: %w", err)
	}

	sb.WriteString(fmt.Sprintf("public class %s\n{\n", tableName))

	for rows.Next() {
		var ordinal int
		var columnName string
		var isNullable string
		var characterMaximumLength sql.NullInt16
		var dataType string
		var columnKey string
		var columnComment sql.NullString

		if err := rows.Scan(
			&ordinal,
			&columnName,
			&isNullable,
			&characterMaximumLength,
			&dataType,
			&columnKey,
			&columnComment,
		); err != nil {
			return "", err
		}

		var cSharpType string

		switch strings.ToLower(dbType) {
		case "mysql":
			cSharpType = mapMySQLToCSharp(dataType)
		case "postgres":
			cSharpType = mapPostgresToCSharp(dataType)
		case "mssql":
			cSharpType = mapMSSQLToCSharp(dataType)
		default:
			cSharpType = "any"
		}

		isNull := strings.EqualFold(isNullable, "YES")

		if opt.Nullable && isNull && isValueType(cSharpType) {
			cSharpType += "?"
		} else if opt.Nullable && isNull && !isValueType(cSharpType) {
			cSharpType += "?"
		}

		fieldName := common.ToPascalCase(columnName)
		if opt.CamelCaseProperties {
			fieldName = common.ToCamelCase(columnName)
		}

		if opt.JsonPropertyName {
			sb.WriteString(fmt.Sprintf(
				"    [JsonPropertyName(\"%s\")]\n",
				columnName,
			))
		}

		getter := ""
		setter := ""
		if opt.Getter {
			getter = "get; "
		}
		if opt.Setter {
			setter = "set; "
		}

		sb.WriteString(fmt.Sprintf(
			"    public %s %s { %s%s}\n",
			cSharpType,
			fieldName,
			getter,
			setter,
		))

		if opt.ExtraSpacing {
			sb.WriteString("\n")
		}
	}

	sb.WriteString("}\n")

	if err := rows.Err(); err != nil {
		return "", err
	}

	return sb.String(), nil
}

func isValueType(csharpType string) bool {
	switch csharpType {
	case "int", "long", "float", "double", "decimal", "bool", "DateTime", "Guid":
		return true
	default:
		return false
	}
}

func mapMySQLToCSharp(dbType string) string {
	switch strings.ToLower(dbType) {
	case "tinyint":
		return "byte"
	case "smallint":
		return "short"
	case "mediumint", "int", "integer":
		return "int"
	case "bigint":
		return "long"

	case "decimal", "numeric":
		return "decimal"
	case "float":
		return "float"
	case "double":
		return "double"

	case "bit", "boolean", "bool":
		return "bool"

	case "date", "datetime", "timestamp":
		return "DateTime"
	case "time":
		return "TimeSpan"

	case "char", "varchar", "text", "tinytext", "mediumtext", "longtext":
		return "string"

	case "json":
		return "string"

	case "binary", "varbinary", "blob", "tinyblob", "mediumblob", "longblob":
		return "byte[]"

	default:
		return "object"
	}
}

func mapMSSQLToCSharp(dbType string) string {
	switch strings.ToLower(dbType) {
	case "bit":
		return "bool"

	case "tinyint":
		return "byte"
	case "smallint":
		return "short"
	case "int":
		return "int"
	case "bigint":
		return "long"

	case "decimal", "numeric", "money", "smallmoney":
		return "decimal"
	case "float":
		return "double"
	case "real":
		return "float"

	case "date", "datetime", "datetime2", "smalldatetime":
		return "DateTime"
	case "time":
		return "TimeSpan"
	case "datetimeoffset":
		return "DateTimeOffset"

	case "char", "nchar", "varchar", "nvarchar", "text", "ntext":
		return "string"

	case "binary", "varbinary", "image":
		return "byte[]"

	case "uniqueidentifier":
		return "Guid"

	default:
		return "object"
	}
}

func mapPostgresToCSharp(dbType string) string {
	switch strings.ToLower(dbType) {
	case "smallint", "int2":
		return "short"
	case "integer", "int4":
		return "int"
	case "bigint", "int8":
		return "long"

	case "numeric", "decimal":
		return "decimal"
	case "real", "float4":
		return "float"
	case "double precision", "float8":
		return "double"

	case "boolean", "bool":
		return "bool"

	case "date", "timestamp", "timestamp without time zone", "timestamp with time zone":
		return "DateTime"
	case "time", "time without time zone":
		return "TimeSpan"

	case "uuid":
		return "Guid"

	case "text", "varchar", "character varying", "char", "character":
		return "string"

	case "json", "jsonb":
		return "string"

	case "bytea":
		return "byte[]"

	default:
		return "object"
	}
}
