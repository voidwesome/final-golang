package api

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// NextDate реализует правила повторения
// now - текущая дата, dstart - дата начала задачи, repeat - правило повторения
// возвращает дату следующего выполнения задачи в формате DateLayout
func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	// проверяем, что правило повторения не пустое
	if strings.TrimSpace(repeat) == "" {
		return "", errors.New("пустое правило повторения")
	}

	// преобразуем дату начала в тип time.Time
	start, err := time.Parse(DateLayout, dstart)
	if err != nil {
		log.Printf("ошибка преобразования даты начала %q: %v\n", dstart, err)
		return "", fmt.Errorf("неверная дата начала: %w", err)
	}

	// разделяем правило повторения по пробелам
	parts := strings.Split(repeat, " ")
	switch parts[0] {
	case "d":
		if len(parts) != 2 {
			return "", errors.New("неверный формат d")
		}
		n, err := strconv.Atoi(parts[1])
		if err != nil || n <= 0 || n > 400 {
			return "", errors.New("неверный интервал d")
		}
		t := start
		for {
			t = t.AddDate(0, 0, n)
			if afterNow(t, now) {
				return t.Format(DateLayout), nil
			}
		}
	case "y":
		if len(parts) != 1 {
			return "", errors.New("неверный формат y")
		}
		t := start
		for {
			t = t.AddDate(1, 0, 0)
			if afterNow(t, now) {
				return t.Format(DateLayout), nil
			}
		}
	default:
		return "", errors.New("неподдерживаемый формат правила повторения")
	}
}

// nextDateHandler обрабатывает HTTP-запрос для вычисления следующей даты задачи
// принимает параметры URL: now (текущая дата), date (дата начала), repeat (правило повторения)
func nextDateHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	nowStr := q.Get("now") // текущая дата, если не указана - берём текущую дату системы
	date := q.Get("date")  // дата начала задачи
	rep := q.Get("repeat") // правило повторения задачи

	var now time.Time
	var err error

	// если параметр now не указан, берем текущую дату системы
	if nowStr == "" {
		now = time.Now()
	} else {
		// преобразуем дату из параметра now
		now, err = time.Parse(DateLayout, nowStr)
		if err != nil {
			writeError(w, fmt.Errorf("неверная дата now: %w", err), http.StatusBadRequest)
			return
		}
	}

	// вычисляем следующую дату выполнения задачи
	next, err := NextDate(now, date, rep)
	if err != nil {
		writeError(w, err, http.StatusBadRequest)
		return
	}

	// отправляем результат как обычный текст
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(next))
}
