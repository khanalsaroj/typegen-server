package query

const (
	CountTablesMSSQL = `
        SELECT COUNT(*)
        FROM sys.tables AS t
        JOIN sys.schemas AS s ON t.schema_id = s.schema_id
        WHERE s.name = @p1;
    `

	DatabaseSizeMSSQL = `
    SELECT
        SUM(a.used_pages) * 8 / 1024 AS SizeMB
    FROM sys.tables t
    JOIN sys.indexes i ON t.object_id = i.object_id
    JOIN sys.partitions p ON i.object_id = p.object_id AND i.index_id = p.index_id
    JOIN sys.allocation_units a ON p.partition_id = a.container_id
    JOIN sys.schemas s ON t.schema_id = s.schema_id
    WHERE s.name = @p1;
    `

	TableNameMSSQL = `
		SELECT
            TABLE_NAME AS "TableName",
            COUNT(*) AS "ColumnCount"
        FROM
            INFORMATION_SCHEMA.COLUMNS
        WHERE
            TABLE_SCHEMA = @p1
        GROUP BY
            TABLE_NAME
    `

	TableColumnDataMSSQL = `
		SELECT
            c.column_id AS ORDINAL_POSITION,
            c.name AS COLUMN_NAME,
            CASE WHEN c.is_nullable = 1 THEN 'YES' ELSE 'NO' END AS IS_NULLABLE,
            CASE
                WHEN t.name IN ('char', 'varchar', 'nchar', 'nvarchar') THEN CAST(c.max_length AS VARCHAR)
                ELSE NULL
            END AS CHARACTER_MAXIMUM_LENGTH,
            t.name AS DATA_TYPE,
            CASE WHEN pk.column_id IS NOT NULL THEN 'PRI' ELSE '' END AS COLUMN_KEY,
            ep.value AS COLUMN_COMMENT
        FROM sys.tables tab
        INNER JOIN sys.columns c ON tab.object_id = c.object_id
        INNER JOIN sys.types t ON c.user_type_id = t.user_type_id
        INNER JOIN sys.schemas s ON tab.schema_id = s.schema_id
        -- Join to find Primary Keys
        LEFT JOIN sys.index_columns pk ON tab.object_id = pk.object_id
            AND c.column_id = pk.column_id
            AND pk.index_id = 1 -- Usually index 1 is the PK
        -- Join to find Extended Properties (Comments)
        LEFT JOIN sys.extended_properties ep ON tab.object_id = ep.major_id
            AND c.column_id = ep.minor_id
            AND ep.name = 'MS_Description'
        WHERE 1=1
          AND s.name = @p1
          AND tab.name = @p2
        ORDER BY c.column_id
          `
)
