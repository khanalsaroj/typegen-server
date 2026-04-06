package golang

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

	structName := req.Prefix + common.ToPascalCase(tbN) + req.Suffix

	var opt domain.GoStructAdvancedOptions
	if err := json.Unmarshal(req.Options, &opt); err != nil {
		return "", fmt.Errorf("invalid Go options: %w", err)
	}

	sb.WriteString(fmt.Sprintf("type %s struct {\n", structName))

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

		goType := mapDBToGoType(dbType, dataType)

		if opt.PointerFields && isNullable == "YES" {
			goType = "*" + goType
		}

		fieldName := common.ToCamelCase(columnName)
		if opt.ExportFields {
			fieldName = common.ToPascalCase(columnName)
		}

		tags := buildTags(columnName, fieldName, opt, isNullable)

		if opt.Comments && columnComment.Valid && strings.TrimSpace(columnComment.String) != "" {
			sb.WriteString(fmt.Sprintf("    // %s\n", columnComment.String))
		}

		sb.WriteString(fmt.Sprintf("    %s %s %s\n", fieldName, goType, tags))

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

func buildTags(columnName, fieldName string, opt domain.GoStructAdvancedOptions, isNullable string) string {
	var tags []string

	if opt.JsonTags {
		jsonTag := columnName
		if opt.OmitEmpty && isNullable == "YES" {
			jsonTag += ",omitempty"
		}
		tags = append(tags, fmt.Sprintf(`json:"%s"`, jsonTag))
	}

	if opt.DBTags {
		tags = append(tags, fmt.Sprintf(`db:"%s"`, columnName))
	}

	if opt.MapstructureTags {
		tags = append(tags, fmt.Sprintf(`mapstructure:"%s"`, columnName))
	}

	if opt.ValidateTags && isNullable == "NO" {
		tags = append(tags, `validate:"required"`)
	}

	if len(tags) == 0 {
		return ""
	}

	return fmt.Sprintf("`%s`", strings.Join(tags, " "))
}

func mapDBToGoType(dbType, dataType string) string {

	db := strings.ToLower(dbType)
	dt := strings.ToLower(dataType)

	switch db {

	case "mysql", "mariadb":
		switch dt {

		case "int", "integer", "mediumint":
			return "int"
		case "bigint":
			return "int64"
		case "smallint":
			return "int16"
		case "tinyint":
			return "int8"

		case "decimal", "numeric":
			return "float64"

		case "float", "double":
			return "float64"

		case "varchar", "text", "char", "longtext", "mediumtext", "tinytext":
			return "string"

		case "date", "datetime", "timestamp":
			return "time.Time"
		case "time":
			return "time.Time"
		case "year":
			return "int"

		case "boolean", "bit":
			return "bool"

		case "json":
			return "map[string]interface{}"

		case "blob", "longblob", "mediumblob", "tinyblob", "binary", "varbinary":
			return "[]byte"

		default:
			return "interface{}"
		}

	case "postgres", "postgresql":
		switch dt {

		case "int", "int4", "serial":
			return "int"
		case "bigint", "int8", "bigserial":
			return "int64"
		case "smallint", "int2":
			return "int16"

		case "numeric", "decimal", "money":
			return "float64"

		case "real", "double precision", "float4", "float8":
			return "float64"

		case "varchar", "text", "char", "character varying":
			return "string"

		case "date":
			return "time.Time"
		case "timestamp", "timestamp without time zone":
			return "time.Time"
		case "timestamptz", "timestamp with time zone":
			return "time.Time"

		case "bool", "boolean":
			return "bool"

		case "uuid":
			return "string"

		case "json", "jsonb":
			return "map[string]interface{}"

		case "bytea":
			return "[]byte"

		case "_int4", "integer[]":
			return "[]int"
		case "_int8", "bigint[]":
			return "[]int64"
		case "_text", "text[]":
			return "[]string"

		default:
			return "interface{}"
		}

	case "mssql", "sqlserver":
		switch dt {

		case "int":
			return "int"
		case "bigint":
			return "int64"
		case "smallint":
			return "int16"
		case "tinyint":
			return "int8"

		case "decimal", "numeric", "money", "smallmoney":
			return "float64"

		case "float", "real":
			return "float64"

		case "varchar", "nvarchar", "char", "nchar", "text", "ntext":
			return "string"

		case "date":
			return "time.Time"
		case "time":
			return "time.Time"
		case "datetime", "datetime2", "smalldatetime":
			return "time.Time"
		case "datetimeoffset":
			return "time.Time"

		case "bit":
			return "bool"

		case "binary", "varbinary", "image", "rowversion":
			return "[]byte"

		case "uniqueidentifier":
			return "string"

		case "xml":
			return "string"

		default:
			return "interface{}"
		}

	default:
		return "interface{}"
	}
}
