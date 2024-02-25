package mssql

import (
	"database/sql"
	"log"

	_ "github.com/denisenkom/go-mssqldb"
)

func InitSQLDB(connectionString string) (*sql.DB, error) {
	db, err := sql.Open("mssql", connectionString)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	log.Println("Connected to SQL Server database")
	return db, nil
}
