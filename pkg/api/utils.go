package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/Westlam/S13/pkg/db"
)

// accepted date format
const dateFormat = "20060102"

// Maximum number of tasks returned
const limit = 50

const (
	errInvalidDateFormat = "invalid format `date`, must be in format %s"
	errInvalidNowFormat  = "invalid format `now`, must be in format %s"
	errInvalidRepeatRule = "invalid RepeatRule"
	// NextDate (Repeat rule BASE):
	errEmptyRepeat        = "repeat is empty"
	errInvalidFormat      = "invalid repeat rule format"
	errInvalidDstart      = "invalid dstart (date) format"
	errInvalidInterval    = "repeat interval is not a valid integer"
	errIntervalOutOfRange = "repeat interval out of range (1..400)"
	// NextDate (Repeat rule *):
)

func writeJson(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(v)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to encode JSON: %v", err), http.StatusInternalServerError)
		return
	}
}

// проверка даты
func checkDate(task *db.Task) error {
	if task == nil {
		return errors.New("task is nil")
	}

	now := time.Now()

	if task.Date == "" {
		task.Date = now.Format(dateFormat)
	}

	t, err := time.Parse(dateFormat, task.Date)
	if err != nil {
		return fmt.Errorf(errInvalidDateFormat, dateFormat)
	}

	// Если дата задачи меньше текущей, обновляем её на сегодня
	if task.Date < now.Format(dateFormat) {
		task.Date = now.Format(dateFormat)
	}

	if task.Repeat != "" {
		if _, err := NextDate(now, task.Date, task.Repeat); err != nil {
			return fmt.Errorf("%s: %w", errInvalidRepeatRule, err)
		}

		if t.Format(dateFormat) < now.Format(dateFormat) {
			nextDate, _ := NextDate(now, task.Date, task.Repeat)
			task.Date = nextDate
		}

		if t.Format(dateFormat) < now.Format(dateFormat) {
			task.Date = now.Format(dateFormat)
		}
	}

	return nil
}
