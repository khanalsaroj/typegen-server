package java

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/khanalsaroj/typegen-server/internal/common"
	"github.com/khanalsaroj/typegen-server/internal/domain"
)

type Record struct{}

func (d *Record) Generate(rows *sql.Rows, req domain.TypeRequest, tbN string, dbType string) (string, error) {
	var sb strings.Builder
	tableName := req.Prefix + common.ToPascalCase(tbN) + req.Suffix

	var opt domain.RecordOptions
	if err := json.Unmarshal(req.Options, &opt); err != nil {
		return "Invalid Record Options", fmt.Errorf("invalid Record options: %w", err)
	}

	if opt.Builder {
		sb.WriteString("@Builder\n")
	}
	sb.WriteString(fmt.Sprintf("public record %s (\n", tableName))
	var fields []string
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

		// Map DB type → Java type
		var javaType string
		if dbType == "mysql" {
			javaType = mapMySQLToJavaType(dataType)
		} else {
			javaType = mapPostgresToJavaType(dataType)
		}

		fieldName := common.ToCamelCase(columnName)

		var fieldSb strings.Builder

		// Swagger annotation
		if columnComment.Valid && strings.TrimSpace(columnComment.String) != "" {
			if opt.SwaggerAnnotations {
				fieldSb.WriteString(fmt.Sprintf(
					"    @Schema(description = \"%s\")\n",
					columnComment.String,
				))
			}
		}

		// Jackson annotation
		if opt.JacksonAnnotations {
			fieldSb.WriteString(fmt.Sprintf(
				"    @JsonProperty(\"%s\")\n",
				fieldName,
			))
		}

		// Field definition (NO comma here)
		fieldSb.WriteString(fmt.Sprintf(
			"    %s %s",
			javaType,
			fieldName,
		))

		fields = append(fields, fieldSb.String())
	}

	if err := rows.Err(); err != nil {
		return "", err
	}

	// Join fields safely (controls comma + spacing)
	separator := ",\n"
	if opt.ExtraSpacing {
		separator = ",\n\n"
	}

	sb.WriteString(strings.Join(fields, separator))
	sb.WriteString("\n) {}")

	if err := rows.Err(); err != nil {
		return "", err
	}

	return sb.String(), nil

}
