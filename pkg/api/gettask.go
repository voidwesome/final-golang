package api

import (
	"net/http"

	"final-golang/pkg/db"
)

// getTaskHandler обрабатывает HTTP-запрос на получение задачи по её идентификатору
// ожидает параметр "id" в URL, возвращает JSON с задачей или ошибкой
func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	// получаем параметр id из строки запроса
	id := r.URL.Query().Get("id")
	if id == "" {
		// если id не указан - возвращаем ошибку в формате JSON
		writeJSON(w, errorResponse{Error: "Не указан идентификатор"})
		return
	}

	// пытаемся получить задачу из базы данных по id
	task, err := db.GetTask(id)
	if err != nil {
		// если задача не найдена - возвращаем ошибку в формате JSON
		writeJSON(w, errorResponse{Error: "Задача не найдена"})
		return
	}

	// отправляем задачу в формате JSON
	writeJSON(w, task)
}
