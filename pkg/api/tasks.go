package api

import (
	"encoding/json"
	"errors"
	"final-golang/pkg/db"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type errorResponse struct {
	Error string `json:"error"`
}

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

type EmptyResponse struct{}

// writeJSON is a helper function to write JSON responses.
func writeJSON(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

// tasksHandler handles GET requests for listing tasks.
func tasksHandler(w http.ResponseWriter, r *http.Request) {
	tasks, err := db.Tasks(50) // Assumes db.DB is used internally by db.Tasks
	if err != nil {
		writeJSON(w, errorResponse{Error: err.Error()})
		return
	}

	writeJSON(w, TasksResp{Tasks: tasks})
}

// getTaskHandler handles GET requests for a single task by ID.
func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSON(w, errorResponse{Error: "Не указан идентификатор"})
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		writeJSON(w, errorResponse{Error: "Задача не найдена"})
		return
	}

	writeJSON(w, task)
}

func postTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeJSON(w, errorResponse{Error: "Ошибка десериализации JSON"})
		return
	}

	if task.Title == "" {
		writeJSON(w, errorResponse{Error: "Не указан заголовок"})
		return
	}

	if task.Date == "today" || task.Date == "" {
		task.Date = time.Now().Format(dateLayout)
	}

	if err := checkDate(&task); err != nil {
		writeJSON(w, errorResponse{Error: err.Error()})
		return
	}

	id, err := db.AddTask(DbConn, &task)
	if err != nil {
		writeJSON(w, errorResponse{Error: err.Error()})
		return
	}

	writeJSON(w, map[string]string{"id": strconv.FormatInt(id, 10)})
}

// updateTaskHandler handles PUT/POST requests to update an existing task.
func UpdateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeJSON(w, errorResponse{Error: "Ошибка десериализации JSON"})
		return
	}

	// Validate mandatory fields for update.
	if task.ID == 0 || task.Title == "" {
		writeJSON(w, errorResponse{Error: "Не указан ID или заголовок"})
		return
	}

	// Handle "today" keyword for the date.
	if task.Date == "today" {
		task.Date = time.Now().Format(dateLayout)
	}

	// Validate the date format and repeat rule.
	if err := checkDate(&task); err != nil {
		writeJSON(w, errorResponse{Error: err.Error()})
		return
	}

	// Update the task in the database.
	if err := db.UpdateTask(&task); err != nil {
		writeJSON(w, errorResponse{Error: err.Error()})
		return
	}

	writeJSON(w, EmptyResponse{})
}

// markDoneHandler handles POST requests to mark a task as done.
// If the task has a repeat rule, its date is updated to the next occurrence.
// Otherwise, the task is deleted.
func MarkDoneHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, errorResponse{Error: "Method not allowed"})
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSON(w, errorResponse{Error: "Не указан ID"})
		return
	}

	t, err := db.GetTask(id)
	if err != nil {
		writeJSON(w, errorResponse{Error: "Задача не найдена"})
		return
	}

	// If no repeat rule, delete the task.
	if t.Repeat == "" {
		if err := db.DeleteTask(id); err != nil {
			writeJSON(w, errorResponse{Error: err.Error()})
			return
		}
		writeJSON(w, struct{}{}) // Empty success response
		return
	}

	// Parse the current task's date to calculate the next occurrence.
	prev, err := time.Parse(dateLayout, t.Date)
	if err != nil {
		writeJSON(w, errorResponse{Error: "Некорректная дата в задаче"})
		return
	}

	// Calculate the next date based on the repeat rule.
	// The 'start' parameter in NextDate is passed as t.Date, which is the current date of the task.
	nextDate, err := NextDate(prev, t.Date, t.Repeat)
	if err != nil {
		writeJSON(w, errorResponse{Error: err.Error()})
		return
	}

	// Update the task's date in the database.
	if err := db.UpdateDate(nextDate, id); err != nil {
		writeJSON(w, errorResponse{Error: err.Error()})
		return
	}

	writeJSON(w, struct{}{}) // Empty success response
}

// deleteTaskHandler handles DELETE requests to delete a task by ID.
func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSON(w, errorResponse{Error: "Не указан ID"})
		return
	}
	if err := db.DeleteTask(id); err != nil {
		writeJSON(w, errorResponse{Error: err.Error()})
		return
	}
	writeJSON(w, struct{}{}) // Empty success response
}

func checkDate(task *db.Task) error {
	now := time.Now().Truncate(24 * time.Hour) // Текущая дата без времени

	// Если дата пуста, устанавливаем её на текущую дату
	if task.Date == "" {
		task.Date = now.Format(dateLayout)
	}

	// Проверяем формат даты (ГГГГММДД)
	parsedDate, err := time.Parse(dateLayout, task.Date)
	if err != nil {
		return errors.New("некорректный формат даты. Ожидается ГГГГММДД")
	}

	// Если дата меньше текущей — корректируем её
	if parsedDate.Before(now) {
		if task.Repeat == "" {
			// Если нет правила повторения — ставим на сегодня
			task.Date = now.Format(dateLayout)
		} else {
			// Есть повторение — вычисляем следующую дату
			nextDate, err := NextDate(now, task.Date, task.Repeat)
			if err != nil {
				return fmt.Errorf("ошибка при вычислении следующей даты для повторения: %w", err)
			}
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
