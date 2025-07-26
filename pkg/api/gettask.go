package api

import (
	"net/http"

	"final-golang/pkg/db"
)

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
