package mssql

import (
	"database/sql"

	_ "github.com/denisenkom/go-mssqldb"
)

// InitSQLDB - Инициализация подключения к базе данных. Реализация для MS SQL Server.
func InitSQLDB(connectionString string) (*sql.DB, error) {
	db, err := sql.Open("mssql", connectionString)
	if err != nil {
		return nil, err
	}

	// Проверка подключения к базе данных
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
