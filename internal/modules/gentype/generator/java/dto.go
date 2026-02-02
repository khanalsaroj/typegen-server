package java

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
	serializable := ""

	tableName := req.Prefix + common.ToPascalCase(tbN) + req.Suffix

	var opt domain.JavaOptions
	if err := json.Unmarshal(req.Options, &opt); err != nil {
		return "Invalid Java Options", fmt.Errorf("invalid Java options: %w", err)
	}

	if opt.Data {
		sb.WriteString("@Data\n")
	} else {
		if opt.Getter {
			sb.WriteString("@Getter\n")
		}
		if opt.Setter {
			sb.WriteString("@Setter\n")
		}
	}
	if opt.NoArgsConstructor {
		sb.WriteString("@NoArgsConstructor\n")
	}

	if opt.AllArgsConstructor {
		sb.WriteString("@AllArgsConstructor\n")
	}

	if opt.Builder {
		sb.WriteString("@Builder\n")
	}

	if opt.Serializable {
		serializable = "implements Serializable"
	}
	sb.WriteString(fmt.Sprintf("public class %s %s{\n", tableName, serializable))

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

		if dbType == "mysql" {
			tsType = mapMySQLToJavaType(dataType)
		} else {
			tsType = mapPostgresToJavaType(dataType)
		}

		fieldName := common.ToCamelCase(columnName)

		if columnComment.Valid && strings.TrimSpace(columnComment.String) != "" {
			if opt.SwaggerAnnotations {
				sb.WriteString(fmt.Sprintf("    @Schema(description = \" %s\")\n", columnComment.String))
			}
		}

		if opt.JacksonAnnotations {
			sb.WriteString(fmt.Sprintf("    @JsonProperty(\"%s\")\n", fieldName))
		}

		sb.WriteString(fmt.Sprintf(
			"    private %s %s;\n",
			tsType,
			fieldName,
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

func mapMySQLToJavaType(mysqlType string) string {
	switch strings.ToLower(mysqlType) {

	case "int", "integer", "mediumint":
		return "Integer"
	case "bigint":
		return "Long"
	case "smallint":
		return "Short"
	case "tinyint":
		return "Byte"
	case "decimal", "numeric":
		return "BigDecimal"
	case "float":
		return "Float"
	case "double":
		return "Double"

	case "varchar", "char", "text", "longtext", "mediumtext", "tinytext":
		return "String"

	case "date":
		return "LocalDate"
	case "datetime", "timestamp":
		return "LocalDateTime"
	case "time":
		return "LocalTime"
	case "year":
		return "Year"

	case "boolean", "bit":
		return "Boolean"

	case "blob", "longblob", "mediumblob", "tinyblob", "binary", "varbinary":
		return "byte[]"

	case "enum":
		return "String"

	default:
		return "Object"
	}
}

func mapPostgresToJavaType(pgType string) string {
	switch strings.ToLower(pgType) {

	case "smallint", "int2":
		return "Short"
	case "integer", "int", "int4", "serial":
		return "Integer"
	case "bigint", "int8", "bigserial":
		return "Long"
	case "decimal", "numeric", "money":
		return "BigDecimal"
	case "real", "float4":
		return "Float"
	case "double precision", "float8":
		return "Double"

	case "varchar", "character varying", "char", "character", "text", "citext":
		return "String"

	case "date":
		return "LocalDate"
	case "time", "time without time zone":
		return "LocalTime"
	case "timestamp", "timestamp without time zone":
		return "LocalDateTime"
	case "timestamptz", "timestamp with time zone":
		return "OffsetDateTime"

	case "boolean", "bool":
		return "Boolean"

	case "bytea":
		return "byte[]"

	case "uuid":
		return "UUID"

	case "json", "jsonb":
		return "String"
	case "xml":
		return "String"

	case "inet", "cidr", "macaddr":
		return "String"

	case "_int4", "integer[]":
		return "Integer[]"
	case "_int8", "bigint[]":
		return "Long[]"
	case "_text", "text[]":
		return "String[]"
	case "_uuid", "uuid[]":
		return "UUID[]"

	case "enum":
		return "String"

	default:
		return "Object"
	}
}
