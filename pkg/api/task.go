package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/Westlam/S13/pkg/db"
)

func taskHandler(w http.ResponseWriter, r *http.Request, database *sql.DB) {
	switch r.Method {
	case http.MethodPost:
		addTaskHandler(w, r, database)
	case http.MethodGet:
		getTaskHandler(w, r, database)
	case http.MethodPut:
		putTaskHandler(w, r, database)
	case http.MethodDelete:
		deleteTaskHandler(w, r, database)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func addTaskHandler(w http.ResponseWriter, r *http.Request, database *sql.DB) {
	var task db.Task

	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeJson(w, map[string]string{"error": "Failed to read request body"})
		return
	}
	defer r.Body.Close()

	if err := json.Unmarshal(body, &task); err != nil {
		writeJson(w, map[string]string{"error": "JSON parsing error"})
		return
	}

	if task.Title == "" {
		writeJson(w, map[string]string{"error": "The title can not be empty"})
		return
	}

	if err := checkDate(&task); err != nil {
		writeJson(w, map[string]string{"error": err.Error()})
		return
	}

	id, err := db.AddTask(&task, database)
	if err != nil {
		writeJson(w, map[string]string{"error": "Task addition failed."})
		return
	}

	task.ID = fmt.Sprintf("%d", id)

	w.WriteHeader(http.StatusCreated)
	writeJson(w, map[string]string{"id": fmt.Sprintf("%d", id)})
}

func getTaskHandler(w http.ResponseWriter, r *http.Request, database *sql.DB) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJson(w, map[string]string{"error": "Id must not be empty"})
		return
	}

	task, err := db.GetTask(id, database)
	if err != nil {
		writeJson(w, map[string]string{"error": err.Error()})
		return
	}
	writeJson(w, task)
}

func putTaskHandler(w http.ResponseWriter, r *http.Request, database *sql.DB) {
	var task db.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeJson(w, map[string]string{"error": "Invalid JSON format"})
		return
	}

	if task.ID == "" {
		writeJson(w, map[string]string{"error": "Id must not be empty"})
		return
	}

	if task.Title == "" {
		writeJson(w, map[string]string{"error": "Title must not be empty"})
		return
	}

	if err := checkDate(&task); err != nil {
		writeJson(w, map[string]string{"error": err.Error()})
		return
	}

	if err := db.UpdateTask(&task, database); err != nil {
		writeJson(w, map[string]string{"error": "Failed to update task"})
		return
	}

	writeJson(w, map[string]interface{}{})
}

func deleteTaskHandler(w http.ResponseWriter, r *http.Request, database *sql.DB) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJson(w, map[string]string{"error": "Id must not be empty"})
		return
	}

	if err := db.DeleteTask(id, database); err != nil {
		writeJson(w, map[string]string{"error": err.Error()})
		return
	}
	writeJson(w, map[string]interface{}{})
}

func taskDoneHandler(w http.ResponseWriter, r *http.Request, database *sql.DB) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJson(w, map[string]string{"error": "Id must not be empty"})
		return
	}

	task, err := db.GetTask(id, database)
	if err != nil {
		writeJson(w, map[string]string{"error": err.Error()})
		return
	}

	if task.Repeat != "" {
		/*
			// Обновляем дату задачи, если она повторяется
			if err := checkDate(task); err != nil {
				writeJson(w, map[string]string{"error": err.Error()})
				return
			}

			// Обновляем задачу в базе данных
			if err := db.UpdateTask(task, database); err != nil {
				writeJson(w, map[string]string{"error": "Failed to update task"})
				return
			}

			// Возвращаем обновленную дату задачи
			writeJson(w, map[string]interface{}{"next_date": task.Date})
		*/
		prev, err := time.Parse(dateFormat, task.Date)
		if err != nil {
			http.Error(w, "Incorrect date format", http.StatusBadRequest)
			return
		}

		next, err := NextDate(prev, task.Date, task.Repeat)
		if err != nil {
			http.Error(w, "Incorrect rule of repeat", http.StatusBadRequest)
			return
		}
		task.Date = next

		if err := db.UpdateTask(task, database); err != nil {
			http.Error(w, "Failed to edit date", http.StatusBadRequest)
			return
		}
	} else {
		// Удаляем задачу, если она не повторяется
		if err := db.DeleteTask(id, database); err != nil {
			writeJson(w, map[string]string{"error": err.Error()})
			return
		}
	}
	writeJson(w, map[string]interface{}{})
}
