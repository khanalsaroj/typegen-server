package java

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/khanalsaroj/typegen-server/internal/common"
	"github.com/khanalsaroj/typegen-server/internal/domain"
	"strings"
)

type XmlAnnotation struct{}

func (d *XmlAnnotation) Generate(rows *sql.Rows, req domain.MapperRequest, tbN string) (string, error) {
	const (
		insertPrefix = "insert_"
		updatePrefix = "update_"
		deletePrefix = "delete_"
	)

	var opt domain.MyBatisOptions
	if err := json.Unmarshal(req.Options, &opt); err != nil {
		return "Invalid Annotation Options", fmt.Errorf("invalid Annotation Options: %w", err)
	}

	skipPrefixes := []string{insertPrefix, updatePrefix, deletePrefix}
	interfaceName := common.ToPascalCase(tbN)
	tableName := tbN

	rowsData, err := scanRows(rows)
	if err != nil {
		return "", fmt.Errorf("failed to scan rows: %w", err)
	}

	var sb strings.Builder
	sb.WriteString("import org.apache.ibatis.annotations.*;\n")
	sb.WriteString("\n")
	sb.WriteString("@Mapper\n")
	sb.WriteString("public interface ")
	sb.WriteString(interfaceName)
	sb.WriteString("Repository {\n")
	generateMyBatisAnnotation(opt, d, &sb, interfaceName, tableName, rowsData, skipPrefixes)
	sb.WriteString("}")

	return sb.String(), nil
}

func generateMyBatisAnnotation(opts domain.MyBatisOptions, d *XmlAnnotation, sb *strings.Builder, interfaceName, tableName string,
	rowsData []domain.SqlData, skipPrefixes []string) {
	if opts.AllCrud {
		d.writeSelectStatement(sb, interfaceName, tableName, rowsData, skipPrefixes)
		d.writeInsertStatement(sb, interfaceName, tableName, rowsData, skipPrefixes)
		d.writeUpdateStatement(sb, interfaceName, tableName, rowsData, skipPrefixes)
		d.writeDeleteStatement(sb, interfaceName, tableName, rowsData)
		return
	}
	if opts.Select {
		d.writeSelectStatement(sb, interfaceName, tableName, rowsData, skipPrefixes)
	}
	if opts.Insert {
		d.writeInsertStatement(sb, interfaceName, tableName, rowsData, skipPrefixes)
	}
	if opts.Update {
		d.writeUpdateStatement(sb, interfaceName, tableName, rowsData, skipPrefixes)
	}
	if opts.Delete {
		d.writeDeleteStatement(sb, interfaceName, tableName, rowsData)
	}
}

func (d *XmlAnnotation) writeSelectStatement(sb *strings.Builder, interfaceName, tableName string,
	rowsData []domain.SqlData, skipPrefixes []string) {
	sb.WriteString("    @Select(\"\"\"\n")
	sb.WriteString("        SELECT\n")
	columnNames := filterColumns(rowsData, skipPrefixes, false)
	writeColumnList(sb, columnNames)

	removeNumberOfLines(sb, 1)

	sb.WriteString(fmt.Sprintf("\n        FROM %s\n", tableName))
	sb.WriteString("        \"\"\")\n")

	sb.WriteString(fmt.Sprintf(
		"    %sResponse select%s(%sResponse param);\n\n",
		interfaceName,
		interfaceName,
		interfaceName,
	))
}

func (d *XmlAnnotation) writeInsertStatement(sb *strings.Builder, interfaceName, tableName string,
	rowsData []domain.SqlData, skipPrefixes []string) {
	sb.WriteString(fmt.Sprintf("    @Insert(\"\"\"\n        INSERT INTO %s (\n", tableName))
	insertColumns := filterColumns(rowsData, skipPrefixes, true)
	writeColumnList(sb, insertColumns)

	sb.WriteString("          insert_ip,\n")
	sb.WriteString("          insert_user_id,\n")
	sb.WriteString("          insert_dtm\n")
	sb.WriteString("        ) VALUES (\n")

	for _, row := range insertColumns {
		sb.WriteString(fmt.Sprintf("            #{%s},\n", common.ToCamelCase(row.ColumnName)))
	}
	sb.WriteString("            #{insertIp},\n")
	sb.WriteString("            #{insertUserId},\n")
	sb.WriteString("            CURRENT_TIMESTAMP(6)\n")
	sb.WriteString("        )\n")
	sb.WriteString("        \"\"\")\n")
	sb.WriteString(fmt.Sprintf("    int insert%s(%sDto dto);\n\n", interfaceName, interfaceName))
}

func (d *XmlAnnotation) writeUpdateStatement(sb *strings.Builder, interfaceName, tableName string,
	rowsData []domain.SqlData, skipPrefixes []string) {
	sb.WriteString(fmt.Sprintf(
		"    @Update(\"\"\"\n        UPDATE %s\n        SET\n",
		tableName,
	))

	updateColumns := filterColumns(rowsData, skipPrefixes, true)
	for _, row := range updateColumns {
		sb.WriteString(fmt.Sprintf(
			"            %s = #{%s},\n",
			row.ColumnName,
			common.ToCamelCase(row.ColumnName),
		))
	}
	sb.WriteString("            update_ip = #{updateIp},\n")
	sb.WriteString("            update_user_id = #{updateUserId},\n")
	sb.WriteString("            update_dtm = CURRENT_TIMESTAMP(6)\n")

	sb.WriteString("        WHERE TRUE\n")

	primaryKeys := getPrimaryKeys(rowsData)
	for _, pk := range primaryKeys {
		sb.WriteString(fmt.Sprintf("            AND %s = #{%s}\n", pk.ColumnName, common.ToCamelCase(pk.ColumnName)))
	}
	sb.WriteString("        \"\"\")\n")
	sb.WriteString(fmt.Sprintf("    int update%s(%sDto dto);\n\n", interfaceName, interfaceName))
}

func (d *XmlAnnotation) writeDeleteStatement(sb *strings.Builder, interfaceName, tableName string,
	rowsData []domain.SqlData) {
	sb.WriteString(fmt.Sprintf("    @Delete(\"\"\"\n        DELETE\n        FROM %s\n        WHERE TRUE\n", tableName))

	primaryKeys := getPrimaryKeys(rowsData)
	for _, pk := range primaryKeys {
		sb.WriteString(fmt.Sprintf("            AND %s = #{%s}\n", pk.ColumnName, common.ToCamelCase(pk.ColumnName)))
	}
	sb.WriteString("        \"\"\")\n")

	sb.WriteString(fmt.Sprintf("    int delete%s(%sDto dto);\n\n", interfaceName))
}
