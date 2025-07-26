package api

import (
	"encoding/json"
	"final-golang/pkg/db"
	"net/http"
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

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	// default 50
	tasks, err := db.Tasks(50)
	if err != nil {
		writeError(w, err, http.StatusInternalServerError)
		return
	}
	writeJSON(w, TasksResp{Tasks: tasks})
}
