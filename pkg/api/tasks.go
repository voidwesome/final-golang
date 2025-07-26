package api

import (
	"net/http"

	"go_final_project/pkg/db"
)

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	tasks, err := db.Tasks(50) // ограничение на 50 задач
	if err != nil {
		writeJSON(w, map[string]any{"error": err.Error()})
		return
	}

	writeJSON(w, TasksResp{Tasks: tasks})
}
