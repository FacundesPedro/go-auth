package domain

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func InitDb(sqlDriver, connectionString string) *sql.DB {
	db, err := sql.Open(sqlDriver, connectionString) //"postgres://username:password@host/dbname")
	if err != nil {
		log.Println(err)
		return nil
	}

	return db
}

func PingDb(sqlDriver string, db *sql.DB) {
	if db == nil {
		log.Fatalf("DB[%s]: %s", sqlDriver, "Could not create a sql.DB instance")
	}
	if err := db.Ping(); err != nil {
		log.Fatalf("DB[%s]: %s", sqlDriver, err)
	} else {
		log.Printf("DB[%s]: Successfully connected!", sqlDriver)
	}
}
