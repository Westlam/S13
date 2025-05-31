package api

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type RepeatRule struct {
	now      time.Time
	start    time.Time
	repeat   string
	interval int // интервал для дней (если применимо)
}

// newRepeatRule создает новый RepeatRule на основе текущей даты, даты начала и правила повторения.
// Возвращает ошибку, если параметры некорректны.
func newRepeatRule(now time.Time, dstart string, repeat string) (*RepeatRule, error) {
	if strings.TrimSpace(repeat) == "" {
		return nil, errors.New(errEmptyRepeat)
	}
	start, err := time.Parse(dateFormat, dstart)
	if err != nil {
		return nil, fmt.Errorf("%s: %v", errInvalidDstart, err)
	}
	return &RepeatRule{
		now:    now,
		start:  start,
		repeat: strings.TrimSpace(repeat),
	}, nil
}

// nextDate вычисляет следующую дату повторения согласно правилу RepeatRule.
// Возвращает строку с датой в формате "YYYYMMDD" или ошибку.
func (r *RepeatRule) nextDate() (string, error) {
	switch {
	case r.repeat == "y":
		next := r.nextYear()
		return next.Format(dateFormat), nil

	case strings.HasPrefix(r.repeat, "w"):
		// TODO -
		return "", errors.New(errInvalidFormat)
	case strings.HasPrefix(r.repeat, "m"):
		// TODO -
		return "", errors.New(errInvalidFormat)

	case r.repeat == "d":
		// формат "d" без числа — ошибка
		return "", errors.New(errInvalidFormat)

	case strings.HasPrefix(r.repeat, "d "):
		intervalStr := strings.TrimSpace(r.repeat[2:])
		interval, err := strconv.Atoi(intervalStr)
		if err != nil {
			return "", errors.New(errInvalidFormat)
		}
		if interval < 1 || interval > 400 {
			return "", errors.New(errIntervalOutOfRange)
		}
		r.interval = interval

		next := r.nextDays()
		return next.Format(dateFormat), nil

	default:
		return "", errors.New(errInvalidFormat)
	}
}

// nextYear вычисляет следующую дату повторения для ежегодного повторения.
// note: Исходя из тестов всегда надо прибавлять хотябы 1 год.
func (r *RepeatRule) nextYear() time.Time {
	next := r.start

	// Если нужно всегда прибавлять хотя бы 1 год, раскомментируйте:
	next = next.AddDate(1, 0, 0)

	for !next.After(r.now) {
		next = next.AddDate(1, 0, 0)
	}
	return next
}

// nextDays вычисляет следующую дату повторения для повторений с интервалом в днях.
// note: Исходя из тестов всегда надо менять дату, даже если она в будущем.
func (r *RepeatRule) nextDays() time.Time {
	next := r.start.AddDate(0, 0, r.interval)
	for !next.After(r.now) {
		next = next.AddDate(0, 0, r.interval)
	}
	return next
}

// NextDate вычисляет дату следующего повторения события.
// Параметры:
// - now: текущая дата (в формате "YYYYMMDD").
// - dstart: дата начала события (в формате "YYYYMMDD").
// - repeat: строка с правилом повторения, например "y" (ежегодно), "d 7" (каждые 7 дней).
// Возвращает строку с датой следующего повторения в формате "YYYYMMDD" или ошибку.
func NextDate(now time.Time, dstart, repeat string) (string, error) {
	rule, err := newRepeatRule(now, dstart, repeat)
	if err != nil {
		return "", err
	}
	return rule.nextDate()
}

// nextDayGetHandler обрабатывает HTTP GET-запросы на получение даты следующего повторения.
// Ожидает параметры "now", "date" и "repeat" в URL-запросе.
// Возвращает дату следующего повторения или ошибку с соответствующим HTTP статусом.
func nextDayGetHandler(w http.ResponseWriter, r *http.Request) {

	//now1, _ := time.Parse(dateFormat, "20250701")
	//dstart1 := ""
	//date1, _ := NextDate(now1, dstart1, "y")
	//fmt.Println(date1)

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	nowStr := r.FormValue("now")
	if nowStr == "" {
		http.Error(w, "Missing parameter 'now'", http.StatusBadRequest)
		return
	}

	now, err := time.Parse(dateFormat, nowStr)
	if err != nil {
		http.Error(w, "Invalid parameter - now", http.StatusBadRequest)
		return
	}

	dstart := r.FormValue("date")
	if dstart == "" {
		http.Error(w, "Invalid parameter - dstart", http.StatusBadRequest)
		return
	}

	repeat := r.FormValue("repeat")
	if repeat == "" {
		http.Error(w, "Invalid parameter - repeat", http.StatusBadRequest)
		return
	}

	date, err := NextDate(now, dstart, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write([]byte(date))
}
