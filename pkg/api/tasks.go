package api

import (
	"encoding/json"
	"final-golang/pkg/db"
	"net/http"
)

// errorResponse описывает структуру ответа с ошибкой
type errorResponse struct {
	Error string `json:"error"`
}

// TasksResp описывает структуру ответа со списком задач
type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

// EmptyResponse — структура для пустого ответа (без данных)
type EmptyResponse struct{}

// writeJSON - вспомогательная функция для отправки ответа в формате JSON.
func writeJSON(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json") // устанавливаем заголовок "Content-Type"
	json.NewEncoder(w).Encode(data)                    // кодируем структуру в JSON и записываем в ответ
}

// tasksHandler обрабатывает HTTP-запрос для получения списка задач
func tasksHandler(w http.ResponseWriter, r *http.Request) {
	// получаем список задач (по умолчанию 50)
	tasks, err := db.Tasks(50)
	if err != nil {
		// в случае ошибки отправляем статус 500 и сообщение об ошибке
		writeError(w, err, http.StatusInternalServerError)
		return
	}
	// отправляем результат в формате JSON
	writeJSON(w, TasksResp{Tasks: tasks})
}
