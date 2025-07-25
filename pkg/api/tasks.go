package api

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/Westlam/S13/pkg/db"
)

type TaskResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func getTasksHandler(w http.ResponseWriter, r *http.Request, database *sql.DB) {

	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	tasks, err := db.GetTasks(limit, database)
	if err != nil {
		writeJson(w, map[string]string{"error": err.Error()})
		return
	}

	today := time.Now().Format(dateFormat)
	for i := range tasks {
		if tasks[i].Date == "" {
			tasks[i].Date = today
		}
	}

	writeJson(w, TaskResp{
		Tasks: tasks,
	})
}
