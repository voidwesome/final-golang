package api

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// afterNow возвращает true, если date > now (по дате без учёта времени)
func afterNow(date, now time.Time) bool {
	d1 := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	d2 := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	return d1.After(d2)
}

// NextDate вычисляет следующую дату задачи с учётом правила repeat
func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	if repeat == "" {
		return "", errors.New("repeat string is empty")
	}

	// Парсим дату dstart
	date, err := time.Parse("20060102", dstart)
	if err != nil {
		return "", fmt.Errorf("invalid dstart date format: %w", err)
	}

	// Разбиваем repeat на части
	parts := strings.Fields(repeat)
	if len(parts) == 0 {
		return "", errors.New("empty repeat")
	}
	rule := parts[0]

	switch rule {
	case "d":
		// Проверяем, что есть второй аргумент - число дней
		if len(parts) < 2 {
			return "", errors.New("repeat 'd' requires a day count")
		}
		days, err := strconv.Atoi(parts[1])
		if err != nil {
			return "", errors.New("invalid day count in repeat 'd'")
		}
		if days < 1 || days > 400 {
			return "", errors.New("day count in 'd' must be between 1 and 400")
		}

		// Прибавляем дни до тех пор, пока date > now
		for {
			date = date.AddDate(0, 0, days)
			if afterNow(date, now) {
				break
			}
		}
		return date.Format("20060102"), nil

	case "y":
		// Прибавляем год пока дата <= now
		for {
			date = date.AddDate(1, 0, 0)
			// Обработка 29 февраля
			if date.Month() == time.February && date.Day() == 29 {
				if !afterNow(date, now) {
					date = date.AddDate(1, 0, 0)
				}
			}
			if afterNow(date, now) {
				break
			}
		}
		return date.Format("20060102"), nil

	case "w":
		if len(parts) < 2 {
			return "", errors.New("repeat 'w' requires week days")
		}
		daysStr := strings.Split(parts[1], ",")
		days := make([]int, 0, len(daysStr))
		for _, dayStr := range daysStr {
			day, err := strconv.Atoi(dayStr)
			if err != nil || day < 1 || day > 7 {
				return "", errors.New("invalid week day in repeat 'w'")
			}
			days = append(days, day)
		}

		// Находим ближайший день недели
		date = findNextWeekday(date, now, days)
		return date.Format("20060102"), nil

	case "m":
		if len(parts) < 2 {
			return "", errors.New("repeat 'm' requires month days")
		}
		daysStr := strings.Split(parts[1], ",")
		days := make([]int, 0, len(daysStr))
		for _, dayStr := range daysStr {
			day, err := strconv.Atoi(dayStr)
			if err != nil {
				return "", errors.New("invalid month day in repeat 'm'")
			}
			days = append(days, day)
		}

		// Находим ближайший день месяца
		date = findNextMonthDay(date, now, days)
		return date.Format("20060102"), nil

	default:
		return "", errors.New("unknown repeat rule: " + rule)
	}
}

// findNextWeekday находит ближайший день недели из списка days (1-7)
func findNextWeekday(date, now time.Time, days []int) time.Time {
	// Приводим дату к началу дня
	date = time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())

	for {
		// Проверяем текущий день недели
		weekday := int(date.Weekday())
		if weekday == 0 { // Воскресенье
			weekday = 7
		}

		// Проверяем, есть ли текущий день в списке дней
		for _, day := range days {
			if weekday == day && afterNow(date, now) {
				return date
			}
		}

		// Переходим к следующему дню
		date = date.AddDate(0, 0, 1)
	}
}

// findNextMonthDay находит ближайший день месяца из списка days
func findNextMonthDay(date, now time.Time, days []int) time.Time {
	// Приводим дату к началу дня
	date = time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())

	for {
		// Проверяем текущий день месяца
		currentDay := date.Day()

		// Проверяем, есть ли текущий день в списке дней
		for _, day := range days {
			if day > 0 {
				// Обычный день месяца
				if currentDay == day && afterNow(date, now) {
					return date
				}
			} else {
				// Отрицательный день - с конца месяца
				lastDay := daysInMonth(date.Year(), int(date.Month()))
				dayFromEnd := lastDay + day + 1
				if currentDay == dayFromEnd && afterNow(date, now) {
					return date
				}
			}
		}

		// Переходим к следующему дню
		date = date.AddDate(0, 0, 1)
	}
}

// daysInMonth возвращает количество дней в месяце
func daysInMonth(year, month int) int {
	return time.Date(year, time.Month(month)+1, 0, 0, 0, 0, 0, time.UTC).Day()
}

func NextDateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	q := r.URL.Query()

	nowStr := q.Get("now")
	if nowStr == "" {
		nowStr = time.Now().Format("20060102")
	}
	now, err := time.Parse("20060102", nowStr)
	if err != nil {
		http.Error(w, "Invalid date format", http.StatusBadRequest)
		return
	}

	dstart := q.Get("date")
	repeat := q.Get("repeat")

	next, err := NextDate(now, dstart, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(next))
}
