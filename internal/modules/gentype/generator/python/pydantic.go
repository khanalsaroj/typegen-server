package python

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/khanalsaroj/typegen-server/internal/common"
	"github.com/khanalsaroj/typegen-server/internal/domain"
)

type PydanticDto struct{}

func (d *PydanticDto) Generate(rows *sql.Rows, req domain.TypeRequest, tbN string, dbType string) (string, error) {

	var sb strings.Builder

	className := req.Prefix + common.ToPascalCase(tbN) + req.Suffix

	var opt domain.PythonPydanticOptions
	if err := json.Unmarshal(req.Options, &opt); err != nil {
		return "Invalid Python Pydantic Options", fmt.Errorf("invalid python pydantic options: %w", err)
	}

	sb.WriteString("from pydantic import BaseModel")

	if opt.Validation {
		sb.WriteString(", Field")
	}
	if opt.StrictTypes {
		sb.WriteString(", StrictInt, StrictStr, StrictBool, StrictFloat")
	}
	sb.WriteString("\n")

	if opt.OptionalFields {
		sb.WriteString("from typing import Optional\n")
	}

	sb.WriteString("\n")

	sb.WriteString(fmt.Sprintf("class %s(BaseModel):\n", className))

	if opt.Docstrings {
		sb.WriteString(fmt.Sprintf("    \"\"\"%s model\"\"\"\n", className))
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
		pyType := mapPydanticType(dbType, dataType, opt.StrictTypes)

		isOpt := opt.OptionalFields || strings.ToLower(isNullable) == "yes"
		if isOpt {
			pyType = fmt.Sprintf("Optional[%s]", pyType)
		}

		if opt.Comments && columnComment.Valid && strings.TrimSpace(columnComment.String) != "" {
			sb.WriteString(fmt.Sprintf("    # %s\n", columnComment.String))
		}

		fieldLine := fmt.Sprintf("    %s: %s", fieldName, pyType)

		var fieldArgs []string

		if opt.Validation {
			if columnComment.Valid && strings.TrimSpace(columnComment.String) != "" {
				fieldArgs = append(fieldArgs, fmt.Sprintf("description=\"%s\"", columnComment.String))
			}
			if opt.AliasGenerator {
				fieldArgs = append(fieldArgs, fmt.Sprintf("alias=\"%s\"", columnName))
			}
		}

		// Default values
		if opt.DefaultValues || isOpt {
			if opt.Validation {
				fieldLine += fmt.Sprintf(" = Field(None%s)", buildFieldArgs(fieldArgs))
			} else {
				fieldLine += " = None"
			}
		} else if len(fieldArgs) > 0 {
			fieldLine += fmt.Sprintf(" = Field(...%s)", buildFieldArgs(fieldArgs))
		}

		sb.WriteString(fieldLine + "\n")

		if opt.ExtraSpacing {
			sb.WriteString("\n")
		}
	}

	if !hasField {
		sb.WriteString("    pass\n")
	}

	if opt.OrmMode ||
		opt.AllowPopulationByFieldName ||
		opt.UseEnumValues ||
		opt.ArbitraryTypesAllowed ||
		opt.AliasGenerator {

		sb.WriteString("\n    class Config:\n")

		if opt.OrmMode {
			sb.WriteString("        orm_mode = True\n")
		}

		if opt.AllowPopulationByFieldName {
			sb.WriteString("        allow_population_by_field_name = True\n")
		}

		if opt.UseEnumValues {
			sb.WriteString("        use_enum_values = True\n")
		}

		if opt.ArbitraryTypesAllowed {
			sb.WriteString("        arbitrary_types_allowed = True\n")
		}

		if opt.AliasGenerator {
			sb.WriteString("        allow_population_by_field_name = True\n")
		}
	}

	if err := rows.Err(); err != nil {
		return "", err
	}

	return sb.String(), nil
}

func mapPydanticType(dbType string, dataType string, strict bool) string {
	switch strings.ToLower(dbType) {
	case "mysql", "mariadb":
		return mapMySQLToPydanticType(dataType, strict)
	case "postgres", "postgresql":
		return mapPostgresToPydanticType(dataType, strict)
	case "mssql", "sqlserver":
		return mapMSSQLToPydanticType(dataType, strict)
	default:
		return "Any"
	}
}

func mapMySQLToPydanticType(t string, strict bool) string {

	switch strings.ToLower(t) {
	case "int", "integer", "mediumint":
		if strict {
			return "StrictInt"
		}
		return "int"

	case "bigint":
		if strict {
			return "StrictInt"
		}
		return "int"

	case "smallint", "tinyint":
		if strict {
			return "StrictInt"
		}
		return "int"

	case "decimal", "numeric":
		return "float" // optionally Decimal

	case "float", "double":
		if strict {
			return "StrictFloat"
		}
		return "float"

	case "varchar", "char", "text", "longtext", "mediumtext", "tinytext":
		if strict {
			return "StrictStr"
		}
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
		if strict {
			return "StrictBool"
		}
		return "bool"

	case "blob", "longblob", "mediumblob", "tinyblob", "binary", "varbinary":
		return "bytes"

	// enum
	case "enum":
		return "str"

	default:
		return "Any"
	}
}

func mapPostgresToPydanticType(t string, strict bool) string {

	switch strings.ToLower(t) {
	case "smallint", "int2":
		if strict {
			return "StrictInt"
		}
		return "int"

	case "integer", "int", "int4", "serial":
		if strict {
			return "StrictInt"
		}
		return "int"

	case "bigint", "int8", "bigserial":
		if strict {
			return "StrictInt"
		}
		return "int"

	case "decimal", "numeric", "money":
		return "float"

	case "real", "float4", "double precision", "float8":
		if strict {
			return "StrictFloat"
		}
		return "float"

	case "varchar", "character varying", "char", "character", "text", "citext":
		if strict {
			return "StrictStr"
		}
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
		if strict {
			return "StrictBool"
		}
		return "bool"

	case "bytea":
		return "bytes"

	case "uuid":
		return "str"

	case "json", "jsonb":
		return "dict"

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

func mapMSSQLToPydanticType(t string, strict bool) string {

	switch strings.ToLower(t) {

	case "int":
		if strict {
			return "StrictInt"
		}
		return "int"

	case "bigint":
		if strict {
			return "StrictInt"
		}
		return "int"

	case "smallint", "tinyint":
		if strict {
			return "StrictInt"
		}
		return "int"

	case "decimal", "numeric", "money", "smallmoney":
		return "float"

	case "float", "real":
		if strict {
			return "StrictFloat"
		}
		return "float"

	case "varchar", "nvarchar", "char", "nchar", "text", "ntext":
		if strict {
			return "StrictStr"
		}
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
		if strict {
			return "StrictBool"
		}
		return "bool"

	case "binary", "varbinary", "image", "rowversion", "timestamp":
		return "bytes"

	case "uniqueidentifier":
		return "str"
	default:
		return "Any"
	}
}

func buildFieldArgs(args []string) string {
	if len(args) == 0 {
		return ""
	}
	return ", " + strings.Join(args, ", ")
}
