package lib

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func InitDB() *sql.DB {
	if os.Getenv("REDIS_PUBSUB") == "true" {
		db, err := sql.Open("sqlite3", "./chatdb.db")
		if err != nil {
			log.Fatal(err)
		}

		sqlStmt := `	
	CREATE TABLE IF NOT EXISTS room (
		name VARCHAR(255) NOT NULL PRIMARY KEY,
		private TINYINT NULL
	);
	`
		_, err = db.Exec(sqlStmt)
		if err != nil {
			log.Fatal("%q: %s\n", err, sqlStmt)
		}

		return db
	}
	return nil
}
