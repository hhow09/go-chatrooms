package lib

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func InitDB() *sql.DB {
	db, err := sql.Open("sqlite3", "./chatdb.db")
	if err != nil {
		log.Fatal(err)
	}

	migration := `	
	CREATE TABLE IF NOT EXISTS room (
		name VARCHAR(255) NOT NULL PRIMARY KEY,
		private TINYINT NULL
	);
	`
	_, err = db.Exec(migration)
	if err != nil {
		log.Fatalf("%q: %s\n", err, migration)
	}
	return db
}
