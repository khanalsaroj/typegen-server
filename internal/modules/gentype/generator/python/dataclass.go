package python

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/khanalsaroj/typegen-server/internal/common"
	"github.com/khanalsaroj/typegen-server/internal/domain"
)

type DataClass struct{}

func (d *DataClass) Generate(rows *sql.Rows, req domain.TypeRequest, tbN string, dbType string) (string, error) {

	var sb strings.Builder

	className := req.Prefix + common.ToPascalCase(tbN) + req.Suffix

	var opt domain.PythonClassOptions
	if err := json.Unmarshal(req.Options, &opt); err != nil {
		return "Invalid Python Class Options", fmt.Errorf("invalid python class options: %w", err)
	}

	if opt.OptionalFields {
		sb.WriteString("from typing import Optional\n\n")
	}

	sb.WriteString(fmt.Sprintf("class %s:\n", className))

	if opt.Docstrings {
		sb.WriteString(fmt.Sprintf("    \"\"\"%s class\"\"\"\n", className))
	}

	type Field struct {
		Name     string
		Type     string
		Comment  string
		Optional bool
	}

	var fields []Field

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

		fieldName := common.ToSnakeCase(columnName)
		pyType := mapDBToPythonType(dbType, dataType)

		isOpt := opt.OptionalFields || strings.ToLower(isNullable) == "yes"
		if isOpt {
			pyType = fmt.Sprintf("Optional[%s]", pyType)
		}

		comment := ""
		if columnComment.Valid {
			comment = columnComment.String
		}

		fields = append(fields, Field{
			Name:     fieldName,
			Type:     pyType,
			Comment:  comment,
			Optional: isOpt,
		})
	}

	if opt.InitMethod {

		sb.WriteString("\n    def __init__(self")

		for _, f := range fields {
			if opt.DefaultValues || f.Optional {
				sb.WriteString(fmt.Sprintf(", %s: %s = None", f.Name, f.Type))
			} else {
				sb.WriteString(fmt.Sprintf(", %s: %s", f.Name, f.Type))
			}
		}
		sb.WriteString("):\n")

		if len(fields) == 0 {
			sb.WriteString("        pass\n")
		} else {
			for _, f := range fields {

				if opt.Comments && strings.TrimSpace(f.Comment) != "" {
					sb.WriteString(fmt.Sprintf("        # %s\n", f.Comment))
				}

				sb.WriteString(fmt.Sprintf("        self.%s = %s\n", f.Name, f.Name))

				if opt.ExtraSpacing {
					sb.WriteString("\n")
				}
			}
		}
	} else {
		for _, f := range fields {

			if opt.Comments && strings.TrimSpace(f.Comment) != "" {
				sb.WriteString(fmt.Sprintf("    # %s\n", f.Comment))
			}

			if opt.DefaultValues || f.Optional {
				sb.WriteString(fmt.Sprintf("    %s: %s = None\n", f.Name, f.Type))
			} else {
				sb.WriteString(fmt.Sprintf("    %s: %s\n", f.Name, f.Type))
			}

			if opt.ExtraSpacing {
				sb.WriteString("\n")
			}
		}
	}

	if opt.ReprMethod && len(fields) > 0 {

		sb.WriteString("\n    def __repr__(self) -> str:\n")

		sb.WriteString("        return f\"")
		sb.WriteString(className + "(")

		for i, f := range fields {
			if i > 0 {
				sb.WriteString(", ")
			}
			sb.WriteString(fmt.Sprintf("%s={self.%s}", f.Name, f.Name))
		}

		sb.WriteString(")\"\n")
	}

	if opt.EqMethod && len(fields) > 0 {

		sb.WriteString("\n    def __eq__(self, other) -> bool:\n")
		sb.WriteString(fmt.Sprintf("        if not isinstance(other, %s):\n", className))
		sb.WriteString("            return False\n")

		sb.WriteString("        return (\n")

		for i, f := range fields {
			if i == len(fields)-1 {
				sb.WriteString(fmt.Sprintf("            self.%s == other.%s\n", f.Name, f.Name))
			} else {
				sb.WriteString(fmt.Sprintf("            self.%s == other.%s and\n", f.Name, f.Name))
			}
		}

		sb.WriteString("        )\n")
	}

	return sb.String(), nil
}
