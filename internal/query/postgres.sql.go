package query

const (
	CountTablesPostgres = `
        SELECT COUNT(*)
		FROM pg_tables
		WHERE schemaname = $1
    `

	DatabaseSizePostgres = `
       SELECT pg_database_size($1) / 1024.0 / 1024.0 AS size_mb;
    `

	TableNamePostgres = `
		SELECT 
			table_name AS "TableName", 
			COUNT(*) AS "ColumnCount"
		FROM 
			information_schema.columns
		WHERE 
			table_schema = $1
		GROUP BY 
			table_name
		ORDER BY 
    		"ColumnCount" DESC;
    `
	ColumnDataPostgres = `
		  SELECT 
			c.ordinal_position as ordinal,
			c.column_name,
			c.is_nullable,
			c.character_maximum_length,
			c.udt_name as dataType,
			CASE 
				WHEN kcu.column_name IS NOT NULL THEN 'PRI'
				ELSE ''
			END AS column_key,
			pgd.description AS column_comment
		FROM information_schema.columns c
		LEFT JOIN information_schema.key_column_usage kcu
			   ON c.table_schema = kcu.table_schema
			  AND c.table_name = kcu.table_name
			  AND c.column_name = kcu.column_name
		LEFT JOIN information_schema.table_constraints tc
			   ON kcu.constraint_name = tc.constraint_name
			  AND tc.constraint_type = 'PRIMARY KEY'
		LEFT JOIN pg_catalog.pg_class pc
			   ON pc.relname = c.table_name
		LEFT JOIN pg_catalog.pg_description pgd
			   ON pgd.objoid = pc.oid
			  AND pgd.objsubid = c.ordinal_position
		WHERE c.table_schema = $1
		  AND c.table_name = $2
		AND c.table_catalog = $3
		ORDER BY c.ordinal_position	  
          `
)
