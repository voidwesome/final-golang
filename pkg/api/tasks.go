package api

import (
	"net/http"

	"final-golang/pkg/db"
)

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
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
