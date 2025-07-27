package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"final-golang/pkg/db"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// DbConn - глобальная переменная для подключения к базе данных
var DbConn *sql.DB

// dateLayout - формат даты для всего приложения (ГГГГММДД)
const dateLayout = "20060102"

// moscowOffset - смещение часового пояса Москвы в секундах (UTC+3)
const moscowOffset = 3 * 60 * 60 // 3 часа в секундах

// addTaskHandler обрабатывает POST-запрос на добавление новой задачи.
func addTaskHandler(w http.ResponseWriter, r *http.Request) {

	// разрешён только метод POST
	if r.Method != http.MethodPost {
		http.Error(w, "Разрешён только POST", http.StatusMethodNotAllowed)
		return
	}

	var task db.Task
	// декодируем JSON из тела запроса в структуру задачи
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		writeJSON(w, errorResponse{Error: "Ошибка десериализации JSON"})
		return
	}

	// проверяем, что заголовок задачи указан
	if task.Title == "" {
		writeJSON(w, errorResponse{Error: "Не указан заголовок задачи"})
		return
	}

	// обработка пустой даты или специального значения "today"
	if task.Date == "" || task.Date == "today" {

		// устанавливаем текущую дату в формате dateLayout
		task.Date = time.Now().Format(dateLayout)
	}

	// проверяем корректность даты и правила повторения
	err = checkDate(&task)
	if err != nil {
		writeJSON(w, errorResponse{Error: err.Error()})
		return
	}

	// добавляем задачу в базу данных (важный шаг!)
	id, err := db.AddTask(&task)
	if err != nil {
		writeJSON(w, errorResponse{Error: "Ошибка добавления в базу"})
		return
	}

	// возвращаем ID новой задачи в JSON-ответе
	writeJSON(w, map[string]string{"id": strconv.FormatInt(id, 10)})
}

// checkDate проверяет и корректирует дату задачи и правило повторения.
func checkDate(task *db.Task) error {
	nowUTC := time.Now().UTC()
	nowMoscow := nowUTC.Add(time.Second * moscowOffset)

	nowMoscow = time.Date(
		nowMoscow.Year(), nowMoscow.Month(), nowMoscow.Day(),
		0, 0, 0, 0, time.FixedZone("MSK", moscowOffset),
	)
	if task.Date == "" {
		task.Date = nowMoscow.Format(dateLayout)
	}

	parsedDate, err := time.Parse(dateLayout, task.Date)
	if err != nil {
		return errors.New("некорректный формат даты. Ожидается ГГГГММДД")
	}

	//если дата задачи раньше текущей - корректируем её
	if parsedDate.Before(nowMoscow) {
		if task.Repeat == "" {
			// если правило повторения отсутствует - ставим дату на сегодня
			task.Date = nowMoscow.Format(dateLayout)
		} else {
			// иначе вычисляем следующую дату согласно правилу повторения
			nextDate, err := NextDate(nowMoscow, task.Date, task.Repeat)
			if err != nil {
				return fmt.Errorf("ошибка при вычислении следующей даты для повторения: %w", err)
			}

			task.Date = nextDate
		}
	}

	// проверяем корректность правила повторения (если оно есть)
	if task.Repeat == "" {
		return nil
	}

	// разбиваем правило повторения на части
	parts := strings.Fields(strings.TrimSpace(task.Repeat))
	if len(parts) == 0 {
		return errors.New("некорректное правило повторения: пусто")
	}

	switch parts[0] {
	case "d": // повторение каждые N дней
		if len(parts) != 2 {
			return errors.New("формат d должен быть: d <число>")
		}
		n, err := strconv.Atoi(parts[1])
		if err != nil || n < 1 || n > 400 {
			return errors.New("некорректный интервал дней для d (ожидается 1-400)")
		}
	case "y": // повторение каждый год
		if len(parts) != 1 {
			return errors.New("формат y не должен содержать параметры")
		}
	case "w": // повторение по дням недели
		if len(parts) != 2 {
			return errors.New("формат w должен быть: w <дни_недели>")
		}
		daysStr := strings.Split(parts[1], ",")
		if len(daysStr) == 0 {
			return errors.New("не указаны дни недели для w")
		}
		for _, day := range daysStr {
			d, err := strconv.Atoi(day)
			if err != nil || d < 1 || d > 7 {
				return errors.New("некорректный день недели для w (ожидается 1-7)")
			}
		}
	default:
		return fmt.Errorf("неподдерживаемый формат повторения: %s", parts[0])
	}

	return nil
}
