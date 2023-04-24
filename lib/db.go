package lib 

import (
	_ "github.com/go-sql-driver/mysql"
	"database/sql"
	"fmt"
	"strings"
)

type TdexStructure struct {
	Endpoint 		*string
	Port 			*string
	User 			*string 
	Pass 			*string 
	TargetDB 		*string
	ExceptTable 	*string
	Object 			*sql.DB
	Schema			map[string][]PerTable
}

type PerTable struct {
	TableName 	string
	General		GeneralInfo
	Columns 	[]ColumnInfo
	Constraints []ConstInfo 	
	Indexes 	[]IndexInfo
	View		ViewInfo
	DDL 		string
}

type GeneralInfo struct {
	TableType 	string
	Engine		sql.NullString
	RowFormat 	sql.NullString
	Collate 	sql.NullString
	Comment 	sql.NullString
}

type ColumnInfo struct {
	ColumnName 		string
	DefaultValue 	sql.NullString
	Nullable		string
	ColumnType		string
	Charset			sql.NullString
	Collation		sql.NullString
	ColumnKey		sql.NullString
	Extra			sql.NullString
	Comment			sql.NullString
}

type ConstInfo struct {
	ConstName		string
	ConstColumn		string
	ConstRefer		string
	ConstDelAction	string
	ConstUpdAction	string
}

type IndexInfo struct {
	IndexName 		string 
	IndexUnique 	int
	IndexCol 		string
}

type ViewInfo struct {
	ViewQuery 	string
	Charset 	string 
	Collate		string
}

func (c TdexStructure) GetDBObject() (*sql.DB, error) {
	DSN := "%s:%s@tcp(%s:%s)/information_schema"
	// Create DB Object
	dbObj, err := sql.Open("mysql",fmt.Sprintf(DSN,*c.User,*c.Pass,*c.Endpoint,*c.Port))
	if err != nil {
		return nil,err
	}

	var result int
	err = dbObj.QueryRow("select 1").Scan(&result)
	if err != nil {
		return nil,err
	}

	return dbObj,nil
}

func (c TdexStructure) GetSchema() (map[string][]PerTable, error) {
	var s map[string][]PerTable
	s = make(map[string][]PerTable)

	getQuery := `
		SELECT
			SCHEMA_NAME 
		FROM
			information_schema.SCHEMATA 
		WHERE
			SCHEMA_NAME NOT IN ( 'information_schema', 'mysql', 'sys', 'performance_schema', 'tmp' ) 
			%s
	`
	if c.TargetDB == nil {
		data, err := c.Object.Query(fmt.Sprintf(getQuery,""))
		if err != nil {
			return s, err
		}
		defer data.Close()

		for data.Next() {
			var schema string
			err := data.Scan(&schema)
			if err != nil {
				return s, err
			}

			s[schema] = []PerTable{}
		}
	} else {
		sa := strings.Split(*c.TargetDB,",")
		var inCon []string
		for _, v := range sa {
			inCon = append(inCon,fmt.Sprintf("'%s'",v))
		}

		data, err := c.Object.Query(fmt.Sprintf(getQuery,fmt.Sprintf("AND SCHEMA_NAME IN (%s)",strings.Join(inCon,","))))
		if err != nil {
			return s, err
		}
		defer data.Close()

		for data.Next() {
			var schema string
			err := data.Scan(&schema)
			if err != nil {
				return s, err
			}

			s[schema] = []PerTable{}
		}

	}

	return s, nil
}

func (c TdexStructure) GetTables(s string) ([]PerTable, error) {
	getQuery := `
		SELECT
			table_name,
			table_type,
			ENGINE,
			row_format,
			table_collation,
			table_comment 
		FROM
			information_schema.TABLES 
		WHERE
			table_schema = ?
			%s
	`

	if c.ExceptTable != nil {
		et := strings.Split(*c.ExceptTable,",")
		var likeCon []string
		for _, v := range et {
			likeCon = append(likeCon,fmt.Sprintf("AND TABLE_NAME NOT LIKE '%s'",v))
		}

		getQuery = fmt.Sprintf(getQuery,strings.Join(likeCon," "))
	} else {
		getQuery = fmt.Sprintf(getQuery,"")
	}
	
	var tables []PerTable
	data, err := c.Object.Query(getQuery,s)
	if err != nil {
		return tables,err
	}
	defer data.Close()

	for data.Next() {
		var tableInfo PerTable
		err := data.Scan(
			&tableInfo.TableName,
			&tableInfo.General.TableType,
			&tableInfo.General.Engine,
			&tableInfo.General.RowFormat,
			&tableInfo.General.Collate,
			&tableInfo.General.Comment,
		)
		if err != nil {
			return tables,err
		}
		tables = append(tables,tableInfo)
	}

	return tables, nil
}

func (p PerTable) GetColumns(o *sql.DB, s string) ([]ColumnInfo, error) {
	getQuery := `
		SELECT
			column_name,
			column_default,
			is_nullable,
			column_type,
			character_set_name,
			collation_name,
			column_key,
			concat(extra," ",if(generation_expression is null,"",generation_expression)),
			column_comment 
		FROM
			information_schema.COLUMNS 
		WHERE
			table_name =? 
			AND table_schema =? 
		ORDER BY
			ordinal_position;
	`

	var columns []ColumnInfo
	data, err := o.Query(getQuery,p.TableName,s)
	if err != nil {
		return columns,err
	}
	defer data.Close()

	for data.Next() {
		var column ColumnInfo
		err := data.Scan(
			&column.ColumnName,
			&column.DefaultValue,
			&column.Nullable,
			&column.ColumnType,
			&column.Charset,
			&column.Collation,
			&column.ColumnKey,
			&column.Extra,
			&column.Comment,
		)
		if err != nil {
			return columns,err
		}

		columns = append(columns,column)
	}

	return columns, nil
}

func (p PerTable) GetIndexes(o *sql.DB, s string) ([]IndexInfo, error) {
	getQuery := `
		SELECT
			INDEX_NAME,	
			NON_UNIQUE,
			GROUP_CONCAT( COLUMN_NAME ORDER BY SEQ_IN_INDEX ASC SEPARATOR ',' ) 
		FROM
			information_schema.STATISTICS 
		WHERE
			TABLE_NAME = ? 
			AND TABLE_SCHEMA = ? 
			AND INDEX_NAME != 'PRIMARY' 
		GROUP BY
			TABLE_SCHEMA,
			TABLE_NAME,
			INDEX_NAME 
		ORDER BY
		INDEX_NAME
	`

	var indexes []IndexInfo
	data, err := o.Query(getQuery,p.TableName,s)
	if err != nil {
		return indexes,err
	}
	defer data.Close()

	for data.Next() {
		var index IndexInfo
		err := data.Scan(
			&index.IndexName,
			&index.IndexUnique,
			&index.IndexCol,
		)
		if err != nil {
			return indexes,err
		}

		indexes = append(indexes,index)
	}

	return indexes, nil
}

func (p PerTable) GetConstraints(o *sql.DB, s string) ([]ConstInfo, error) {
	getQuery := `
		SELECT
			x.constraint_name AS constraint_key,
			group_concat( x.column_name ) AS con_column,
			concat( x.referenced_table_name, '.', x.referenced_column_name ) AS refer_info,
			y.DELETE_RULE,
			y.UPDATE_RULE 
		FROM
			information_schema.KEY_COLUMN_USAGE x
			INNER JOIN information_schema.REFERENTIAL_CONSTRAINTS y ON x.constraint_name = y.constraint_name 
		WHERE
			x.table_name = ?
			AND x.CONSTRAINT_SCHEMA = ? 
			AND x.constraint_name <> 'PRIMARY' 
		GROUP BY
			x.constraint_name
	`

	var consts []ConstInfo
	data, err := o.Query(getQuery,p.TableName,s)
	if err != nil {
		return consts,err
	}
	defer data.Close()

	for data.Next() {
		var constraint ConstInfo
		err := data.Scan(
			&constraint.ConstName,
			&constraint.ConstColumn,
			&constraint.ConstRefer,
			&constraint.ConstDelAction,
			&constraint.ConstUpdAction,
		)
		if err != nil {
			return consts,err
		}

		consts = append(consts,constraint)
	}

	return consts, nil
}

func (p PerTable) GetViewSQL(o *sql.DB, s string) (ViewInfo, error) {
	getQuery := `
		SHOW CREATE TABLE %s.%s 
	`

	var v ViewInfo
	var viewName string
	err := o.QueryRow(fmt.Sprintf(getQuery,s,p.TableName)).Scan(&viewName,&v.ViewQuery,&v.Charset,&v.Collate)
	if err != nil {
		return v, err
	}

	return v, nil
}

func (p PerTable) GetTableDDL(o *sql.DB,s string) (string,error) {
	getQuery := `
		SHOW CREATE TABLE %s.%s
	`

	var ddl string 
	var tableName string
	err := o.QueryRow(fmt.Sprintf(getQuery,s,p.TableName)).Scan(&tableName,&ddl)
	if err != nil {
		return ddl, err
	}

	return ddl, nil
}