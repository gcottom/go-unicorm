package unicorm

import "database/sql"

const STATEMENT = "Statement"

type structInfo struct {
	FieldType string
	FieldName string
}

type Table[T any] struct {
	TableName  string
	RepoStruct T
	DBClient   *sql.DB
	structInfo []structInfo
}
