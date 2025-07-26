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

var DbConn *sql.DB

const dateLayout = "20060102"

const moscowOffset = 3 * 60 * 60 // 3 часа в секундах

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	var task db.Task
	err := json.NewDecoder(r.Body).Decode(&task)
	fmt.Printf("декодирование %v\n", task)
	if err != nil {
		writeJSON(w, errorResponse{Error: "Ошибка десериализации JSON"})
		return
	}

	// Validate task title (mandatory field)
	if task.Title == "" {
		writeJSON(w, errorResponse{Error: "Не указан заголовок задачи"})
		return
	}

	// Handle empty date or "today"
	if task.Date == "" || task.Date == "today" {
		task.Date = time.Now().Format(dateLayout)
	}
	fmt.Printf("измененная дата таски - %v\n", task)

	// Validate date and repeat rule
	err = checkDate(&task)
	if err != nil {
		writeJSON(w, errorResponse{Error: err.Error()})
		return
	}
	fmt.Printf("форматрирование - %v\n", task)

	// Add the task to database (важно!)
	id, err := db.AddTask(&task)
	if err != nil {
		writeJSON(w, errorResponse{Error: "Ошибка добавления в базу"})
		return
	}

	writeJSON(w, map[string]string{"id": strconv.FormatInt(id, 10)})
}

func checkDate(task *db.Task) error {
	nowUTC := time.Now().UTC()

	// Добавляем сдвиг
	nowMoscow := nowUTC.Add(time.Second * moscowOffset)

	// Обнуляем время по московскому времени
	nowMoscow = time.Date(
		nowMoscow.Year(), nowMoscow.Month(), nowMoscow.Day(),
		0, 0, 0, 0, time.FixedZone("MSK", moscowOffset),
	)

	// Если дата пуста, устанавливаем её на текущую дату
	if task.Date == "" {
		task.Date = nowMoscow.Format(dateLayout)
	}

	// Проверяем формат даты (ГГГГММДД)
	parsedDate, err := time.Parse(dateLayout, task.Date)
	if err != nil {
		return errors.New("некорректный формат даты. Ожидается ГГГГММДД")
	}

	// Если дата меньше текущей — корректируем её
	if parsedDate.Before(nowMoscow) {
		if task.Repeat == "" {
			fmt.Printf("дата если правило - пустой %s\n", task.Date)
			// Если нет правила повторения — ставим на сегодня
			task.Date = nowMoscow.Format(dateLayout)
		} else {
			fmt.Printf("дата иначе - %s\n", task.Date)
			// Есть повторение — вычисляем следующую дату
			nextDate, err := NextDate(nowMoscow, task.Date, task.Repeat)
			if err != nil {
				return fmt.Errorf("ошибка при вычислении следующей даты для повторения: %w", err)
			}
			fmt.Printf("следующая дата %s\n", nextDate)

			task.Date = nextDate
		}
	}

	// Проверка правила повторения (если оно есть)
	if task.Repeat == "" {
		return nil
	}

	parts := strings.Fields(strings.TrimSpace(task.Repeat))
	if len(parts) == 0 {
		return errors.New("некорректное правило повторения: пусто")
	}

	switch parts[0] {
	case "d":
		if len(parts) != 2 {
			return errors.New("формат d должен быть: d <число>")
		}
		n, err := strconv.Atoi(parts[1])
		if err != nil || n < 1 || n > 400 {
			return errors.New("некорректный интервал дней для d (ожидается 1-400)")
		}
	case "y":
		if len(parts) != 1 {
			return errors.New("формат y не должен содержать параметры")
		}
	case "w":
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
