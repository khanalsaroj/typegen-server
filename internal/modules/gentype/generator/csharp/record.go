package csharp

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

	var opt domain.CSharpRecordOptions
	if err := json.Unmarshal(req.Options, &opt); err != nil {
		return "Invalid CSharp Record Options", fmt.Errorf("invalid CSharp record options: %w", err)
	}

	type field struct {
		Name       string
		DbName     string
		CSharpType string
		IsNullable bool
	}

	var fields []field

	for rows.Next() {
		var ordinal int
		var columnName string
		var isNullable string
		var charMax sql.NullInt16
		var dataType string
		var columnKey string
		var columnComment sql.NullString

		if err := rows.Scan(
			&ordinal,
			&columnName,
			&isNullable,
			&charMax,
			&dataType,
			&columnKey,
			&columnComment,
		); err != nil {
			return "", err
		}

		var csharpType string
		switch strings.ToLower(dbType) {
		case "mysql":
			csharpType = mapMySQLToCSharp(dataType)
		case "postgres", "postgresql":
			csharpType = mapPostgresToCSharp(dataType)
		case "mssql", "sqlserver":
			csharpType = mapMSSQLToCSharp(dataType)
		default:
			csharpType = "object"
		}

		isNull := strings.EqualFold(isNullable, "YES")

		if opt.Nullable && isNull {
			csharpType = makeNullableCSharpType(csharpType)
		}

		propName := common.ToPascalCase(columnName)
		if opt.CamelCaseProperties {
			propName = common.ToCamelCase(columnName)
		}

		fields = append(fields, field{
			Name:       propName,
			DbName:     columnName,
			CSharpType: csharpType,
			IsNullable: isNull,
		})
	}

	if opt.Positional {
		sb.WriteString(fmt.Sprintf("public record %s(\n", tableName))
		for i, f := range fields {
			sb.WriteString(fmt.Sprintf("    %s %s", f.CSharpType, f.Name))
			if i < len(fields)-1 {
				sb.WriteString(",")
			}
			sb.WriteString("\n")
		}
		sb.WriteString(");\n")
		return sb.String(), nil
	}

	if opt.JsonPropertyName {
		sb.WriteString("using System.Text.Json.Serialization;\n\n")
	}

	sb.WriteString(fmt.Sprintf("public record %s\n{\n", tableName))

	for _, f := range fields {
		if opt.JsonPropertyName {
			sb.WriteString(fmt.Sprintf(
				"    [JsonPropertyName(\"%s\")]\n",
				f.DbName,
			))
		}

		accessor := "init;"
		if !opt.WithInit {
			accessor = "set;"
		}

		sb.WriteString(fmt.Sprintf(
			"    public %s %s { get; %s }\n",
			f.CSharpType,
			f.Name,
			accessor,
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

func makeNullableCSharpType(t string) string {
	switch t {
	case "int", "long", "short", "byte", "float", "double", "decimal", "bool", "DateTime", "Guid":
		return t + "?"
	default:
		if strings.HasSuffix(t, "?") {
			return t
		}
		return t + "?"
	}
}
