package db

import (
	"database/sql"
	"log"
	"os"

	"github.com/Westlam/S13/config"
	_ "modernc.org/sqlite"
)

func Init(config config.Config) (*sql.DB, error) {
	dbFile := config.TODO_DBFILE
	// Проверяем, существует ли база данных
	_, err := os.Stat(dbFile)
	install := os.IsNotExist(err)

	log.Println("Using database file:", dbFile)

	db, err := sql.Open("sqlite", dbFile)
	if err != nil {
		return nil, err
	}

	if install {
		query := `
		CREATE TABLE IF NOT EXISTS scheduler (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			date TEXT NOT NULL,
			title TEXT NOT NULL,
			comment TEXT,
			repeat TEXT CHECK(length(repeat) <= 128)
		);
		CREATE INDEX IF NOT EXISTS idx_date ON scheduler(date);
		`
		_, err := db.Exec(query)
		if err != nil {
			return nil, err
		}
		log.Println("Database created")
	}

	return db, nil
}
