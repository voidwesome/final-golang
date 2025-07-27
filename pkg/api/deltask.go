package api

import (
	"net/http"

	"final-golang/pkg/db"
)

// delTaskHandler обрабатывает DELETE-запросы на удаление задачи по её ID.
func delTaskHandler(w http.ResponseWriter, r *http.Request) {
	// получаем ID задачи из параметров запроса
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSON(w, errorResponse{Error: "Не указан ID"})
		return
	}

	// удаляем задачу из базы данных
	if err := db.DeleteTask(id); err != nil {
		// в случае ошибки возвращаем её в формате JSON
		writeJSON(w, errorResponse{Error: err.Error()})
		return
	}

	// возвращаем пустой JSON в случае успешного удаления
	writeJSON(w, struct{}{})
}
