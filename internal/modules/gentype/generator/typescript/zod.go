package typescript

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/khanalsaroj/typegen-server/internal/common"
	"github.com/khanalsaroj/typegen-server/internal/domain"
)

type Zod struct{}

func (z Zod) Generate(rows *sql.Rows, req domain.TypeRequest, tbN string, dbType string) (string, error) {

	var sb strings.Builder

	tableName := req.Prefix + common.ToPascalCase(tbN) + req.Suffix

	var opt domain.ZodOptions
	if err := json.Unmarshal(req.Options, &opt); err != nil {
		return "Invalid Zod Options", fmt.Errorf("invalid Zod options: %w", err)
	}

	sb.WriteString("export const ")
	sb.WriteString(tableName)
	sb.WriteString("Schema = z.object({\n")

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

		var zodType string
		if dbType == "mysql" {
			zodType = mapMySQLToZod(dataType)
		} else {
			zodType = mapPostgresToZod(dataType)
		}

		fieldName := common.ToCamelCase(columnName)

		sb.WriteString("  ")
		sb.WriteString(fieldName)
		sb.WriteString(": z.")
		sb.WriteString(zodType)

		// string trimming
		if opt.Trim && zodType == "string()" {
			sb.WriteString(".trim()")
		}

		// max length
		if opt.MaxValue &&
			characterMaximumLength.Valid &&
			zodType == "string()" {
			sb.WriteString(fmt.Sprintf(".max(%d)", characterMaximumLength.Int16))
		}

		// nullability handling (ORDER MATTERS)
		if opt.Nullish {
			sb.WriteString(".nullish()")
		} else {
			if opt.Nullable || isNullable == "YES" {
				sb.WriteString(".nullable()")
			}
			if opt.AllOptional {
				sb.WriteString(".optional()")
			}
		}

		sb.WriteString(",")

		// comments
		if opt.Comments && columnComment.Valid && strings.TrimSpace(columnComment.String) != "" {
			sb.WriteString(fmt.Sprintf(" // %s", columnComment.String))
		}

		sb.WriteString("\n")
	}

	sb.WriteString("}).strict();\n")

	if opt.ExportAllTypes {
		sb.WriteString("\nexport type ")
		sb.WriteString(tableName)
		sb.WriteString(" = z.infer<typeof ")
		sb.WriteString(tableName)
		sb.WriteString("Schema>;\n")
	}

	return sb.String(), nil

}

func mapMySQLToZod(mysqlType string) string {
	switch strings.ToLower(mysqlType) {

	case "int", "bigint", "smallint", "tinyint", "decimal", "float", "double", "long":
		return "number()"

	case "varchar", "char", "text", "longtext":
		return "string()"

	case "datetime", "timestamp", "date":
		return "date()"

	case "boolean", "bit":
		return "bool()"

	case "enum":
		return "enum([])"

	default:
		return "string()"
	}
}

func mapPostgresToZod(pgType string) string {
	switch strings.ToLower(pgType) {

	case "smallint", "int2",
		"integer", "int", "int4", "serial",
		"bigint", "int8", "bigserial",
		"decimal", "numeric", "money",
		"real", "float4",
		"double precision", "float8":
		return "number()"

	case "varchar", "character varying",
		"char", "character",
		"text", "citext",
		"inet", "cidr", "macaddr",
		"xml":
		return "string()"

	case "boolean", "bool":
		return "bool()"

	case "date",
		"timestamp", "timestamp without time zone",
		"timestamptz", "timestamp with time zone",
		"time", "time without time zone":
		return "date()"

	case "uuid":
		return "string().uuid()"

	case "json", "jsonb":
		return "any()"

	case "bytea":
		return "instanceof(Uint8Array)"

	case "integer[]", "_int4":
		return "array(z.number())"
	case "bigint[]", "_int8":
		return "array(z.number())"
	case "text[]", "_text":
		return "array(z.string())"
	case "uuid[]", "_uuid":
		return "array(z.string().uuid())"

	case "enum":
		return "enum([])"

	default:
		return "string()"
	}
}
