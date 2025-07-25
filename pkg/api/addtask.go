package api

import (
	"database/sql"
	"encoding/json"
	"final-golang/pkg/db"
	"net/http"
	"strconv"
	"time"
)

var DbConn *sql.DB

const dateLayout = "20060102"

func AddTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	var task db.Task
	err := json.NewDecoder(r.Body).Decode(&task)
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

	// Validate date and repeat rule
	err = checkDate(&task)
	if err != nil {
		writeJSON(w, errorResponse{Error: err.Error()})
		return
	}

	// Add the task to database (важно!)
	id, err := db.AddTask(db.DB, &task)
	if err != nil {
		writeJSON(w, errorResponse{Error: "Ошибка добавления в базу"})
		return
	}

	writeJSON(w, map[string]string{"id": strconv.FormatInt(id, 10)})
}
