package java

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/khanalsaroj/typegen-server/internal/common"
	"github.com/khanalsaroj/typegen-server/internal/domain"
)

type Xml struct{}

func (d *Xml) Generate(rows *sql.Rows, req domain.MapperRequest, tbN string) (string, error) {
	const (
		insertPrefix = "insert_"
		updatePrefix = "update_"
		deletePrefix = "delete_"
	)

	var opt domain.MyBatisOptions
	if err := json.Unmarshal(req.Options, &opt); err != nil {
		return "Invalid MyBatis Options", fmt.Errorf("invalid MyBatis Options: %w", err)
	}

	skipPrefixes := []string{insertPrefix, updatePrefix, deletePrefix}
	interfaceName := common.ToPascalCase(tbN)
	tableName := tbN

	rowsData, err := scanRows(rows)
	if err != nil {
		return "", fmt.Errorf("failed to scan rows: %w", err)
	}

	var sb strings.Builder
	sb.WriteString("<mapper namespace=\"")
	sb.WriteString(interfaceName)
	sb.WriteString("Repository\">\n")

	generateMyBatis(opt, d, &sb, interfaceName, tableName, rowsData, skipPrefixes)
	sb.WriteString("</mapper>")

	return sb.String(), nil
}

func generateMyBatis(opts domain.MyBatisOptions, d *Xml, sb *strings.Builder, interfaceName, tableName string,
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

func (d *Xml) writeSelectStatement(sb *strings.Builder, interfaceName, tableName string,
	rowsData []domain.SqlData, skipPrefixes []string) {
	sb.WriteString(fmt.Sprintf(`    <select id="select%s" parameterType="%sResponse">`,
		interfaceName, interfaceName))
	sb.WriteString("\n        SELECT\n")

	columnNames := filterColumns(rowsData, skipPrefixes, false)
	writeColumnList(sb, columnNames)

	removeNumberOfLines(sb, 1)

	sb.WriteString(fmt.Sprintf("\n        FROM %s", tableName))
	sb.WriteString("\n    </select>\n\n")
}

func (d *Xml) writeInsertStatement(sb *strings.Builder, interfaceName, tableName string,
	rowsData []domain.SqlData, skipPrefixes []string) {
	sb.WriteString(fmt.Sprintf(`    <insert id="insert%s" parameterType="%sDto">`,
		interfaceName, interfaceName))
	sb.WriteString(fmt.Sprintf("\n        INSERT INTO %s (", tableName))
	sb.WriteString("\n")

	insertColumns := filterColumns(rowsData, skipPrefixes, true)
	writeColumnList(sb, insertColumns)

	sb.WriteString("          insert_ip,\n")
	sb.WriteString("          insert_user_id,\n")
	sb.WriteString("          insert_dtm)\n")
	sb.WriteString("        VALUES (\n")

	for _, row := range insertColumns {
		sb.WriteString(fmt.Sprintf("          #{%s},\n",
			common.ToCamelCase(row.ColumnName)))
	}

	sb.WriteString("          #{insertIp},\n")
	sb.WriteString("          #{insertUserId},\n")
	sb.WriteString("          CURRENT_TIMESTAMP(6));\n")
	sb.WriteString("    </insert>\n\n")
}

func (d *Xml) writeUpdateStatement(sb *strings.Builder, interfaceName, tableName string,
	rowsData []domain.SqlData, skipPrefixes []string) {
	sb.WriteString(fmt.Sprintf(`    <update id="update%s">`, interfaceName))
	sb.WriteString(fmt.Sprintf("\n        UPDATE %s", tableName))
	sb.WriteString("\n        SET\n")

	updateColumns := filterColumns(rowsData, skipPrefixes, true)
	for i, row := range updateColumns {
		sb.WriteString(fmt.Sprintf("          %s= #{%s}",
			row.ColumnName, common.ToCamelCase(row.ColumnName)))
		if i < len(updateColumns)-1 {
			sb.WriteString(",")
		}
		sb.WriteString("\n")
	}

	sb.WriteString("          update_ip = #{updateIp},\n")
	sb.WriteString("          update_user_id = #{updateUserId},\n")
	sb.WriteString("          update_dtm = CURRENT_TIMESTAMP(6)\n")
	sb.WriteString("        WHERE TRUE\n")

	primaryKeys := getPrimaryKeys(rowsData)
	for _, pk := range primaryKeys {
		sb.WriteString(fmt.Sprintf("            AND %s = #{%s}\n",
			pk.ColumnName, common.ToCamelCase(pk.ColumnName)))
	}

	sb.WriteString("    </update>\n\n")
}

func (d *Xml) writeDeleteStatement(sb *strings.Builder, interfaceName, tableName string,
	rowsData []domain.SqlData) {
	sb.WriteString(fmt.Sprintf(`    <delete id="delete%s">`, interfaceName))
	sb.WriteString("\n        DELETE")
	sb.WriteString(fmt.Sprintf("\n        FROM %s", tableName))
	sb.WriteString("\n        WHERE TRUE\n")

	primaryKeys := getPrimaryKeys(rowsData)
	for _, pk := range primaryKeys {
		sb.WriteString(fmt.Sprintf("            AND %s = #{%s}\n",
			pk.ColumnName, common.ToCamelCase(pk.ColumnName)))
	}

	sb.WriteString("    </delete>\n")
}
