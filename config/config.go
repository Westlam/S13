package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

// Config хранит конфигурацию приложения.
type Config struct {
	TODO_PORT     string
	TODO_DBFILE   string
	//WEB_DIR       string
	//TODO_PASSWORD string
}

// NewConfig создает новую конфигурацию, считывая переменные окружения.
func New() Config {
	return Config{
		TODO_PORT:   getPort(),
		TODO_DBFILE: getDatabaseFilePath(),
	}
}

// getPort возвращает порт, на котором должен работать сервер.
// Если переменная окружения TODO_PORT не установлена, возвращает значение по умолчанию.
func getPort() string {
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = "7540" // значение по умолчанию
	}
	return port
}

// getDatabaseFilePath возвращает путь к файлу базы данных.
// Если переменная окружения TODO_DBFILE установлена и не пустая,
// используется её значение. В противном случае используется значение по умолчанию.
func getDatabaseFilePath() string {
	dbFile := os.Getenv("TODO_DBFILE")

	if dbFile == "" {
		cwd, err := os.Getwd()
		if err != nil {
			log.Fatal(fmt.Errorf("error getting current working directory: %w", err))
		}
		dbFile = filepath.Join(cwd, "scheduler.db")
	}
	
	return dbFile
}
