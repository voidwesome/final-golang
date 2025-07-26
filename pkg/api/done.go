package api

import (
	"net/http"
	"time"

	"final-golang/pkg/db"
)

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
