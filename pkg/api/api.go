package api

import (
	"database/sql"
	"net/http"
)

func Init(db *sql.DB) {
	http.HandleFunc("/api/nextdate", nextDayGetHandler)

	http.HandleFunc("/api/task", func(w http.ResponseWriter, r *http.Request) {
		taskHandler(w, r, db)
	})
	http.HandleFunc("/api/tasks", func(w http.ResponseWriter, r *http.Request) {
		getTasksHandler(w, r, db)
	})
	http.HandleFunc("/api/task/done", func(w http.ResponseWriter, r *http.Request) {
		taskDoneHandler(w, r, db)
	})
}
