package api

import (
	"net/http"
	"time"
)

// DateLayout - формат даты для всего приложения (YYYYMMDD)
const DateLayout = "20060102"

// writeError отправляет HTTP-ответ с ошибкой в формате JSON и указанным кодом состояния
func writeError(w http.ResponseWriter, err error, code int) {
	w.WriteHeader(code)
	writeJSON(w, map[string]string{"error": err.Error()})
}

// afterNow сравнивает две даты без учета времени суток и проверяет,
// находится ли первая дата (date) строго после текущей (now)
func afterNow(date, now time.Time) bool {
	// создаем объекты времени только с датой (обнуляем часы, минуты, секунды и наносекунды)
	d := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.Local)
	n := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)

	// возвращает true, если date > now (только по дате, время игнорируется)
	return d.After(n)
}
