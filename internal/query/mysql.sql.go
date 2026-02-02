package query

const (
	CountTablesMySQL = `
        SELECT COUNT(*)
        FROM information_schema.tables
        WHERE table_schema = ?
    `

	DatabaseSizeMySQL = `
        SELECT SUM(data_length + index_length) / 1024 / 1024
        FROM information_schema.tables
        WHERE table_schema = ?
    `

	TableNameMySQL = `
		SELECT 
			table_name AS "TableName", 
			COUNT(*) AS "ColumnCount"
		FROM 
			information_schema.columns
		WHERE 
			table_schema = ?
		GROUP BY 
			table_name
    `

	TableColumnDataMySQL = `
		SELECT ORDINAL_POSITION, 
		       COLUMN_NAME,
		        IS_NULLABLE,
		        CHARACTER_MAXIMUM_LENGTH,
		        DATA_TYPE, 
		        COLUMN_KEY, 
		        COLUMN_COMMENT
		 FROM INFORMATION_SCHEMA.COLUMNS 
		 where table_schema = ?
			 and table_name = ?
			 order by ORDINAL_POSITION
          `
)
