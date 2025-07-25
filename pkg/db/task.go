package db

import (
	"database/sql"
	"fmt"
)

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment,omitempty"`
	Repeat  string `json:"repeat,omitempty"`
}

func AddTask(task *Task, db *sql.DB) (int64, error) {
	var id int64
	query := "INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)"
	res, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err == nil {
		id, err = res.LastInsertId()
	}

	return id, err
}

func GetTasks(limit int, db *sql.DB) ([]*Task, error) {
	rows, err := db.Query(`SELECT id, date, title, comment, repeat
	FROM scheduler
	ORDER BY date ASC
	LIMIT ?`, limit)
	if err != nil {
		return nil, fmt.Errorf("error while SELECT %w", err)
	}
	defer rows.Close()

	var tasks []*Task
	for rows.Next() {
		var t Task
		if err := rows.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat); err != nil {
			return nil, fmt.Errorf("error searching for nearby tasks: %w", err)
		}
		tasks = append(tasks, &t)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error searching for nearby tasks: %w", err)
	}

	if tasks == nil {
		tasks = make([]*Task, 0)
	}
	return tasks, nil
}

func GetTask(id string, db *sql.DB) (*Task, error) {
	var task Task
	query := `SELECT id, date, title, comment, repeat
              FROM scheduler
              WHERE id = ?`
	err := db.QueryRow(query, id).Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("task not found")
		}
		return nil, fmt.Errorf("failed to get task: %w", err)
	}

	return &task, nil
}

func UpdateTask(task *Task, db *sql.DB) error {
	// параметры пропущены, не забудьте указать WHERE
	query := `UPDATE scheduler SET date=?, title=?, comment=?, repeat=? WHERE id=?`
	res, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		return err
	}
	// метод RowsAffected() возвращает количество записей к которым
	// был применена SQL команда
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf(`incorrect id for updating task`)
	}
	return nil
}

func DeleteTask(id string, db *sql.DB) error {
	if id == "" {
		return fmt.Errorf("Id must not be empty")
	}
	res, err := db.Exec(`DELETE FROM scheduler WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("Delete task: %w", err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("Counting deleted rows: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("Task not found: %s", id)
	}
	return nil
}
