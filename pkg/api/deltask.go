package api

import (
	"net/http"

	"final-golang/pkg/db"
)

// deleteTaskHandler handles DELETE requests to delete a task by ID.
func delTaskHandler(w http.ResponseWriter, r *http.Request) {
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
