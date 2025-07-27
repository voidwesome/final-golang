package api

import "net/http"

// taskHandler обрабатывает HTTP-запросы к одному ресурсу "task"
// в зависимости от метода HTTP (POST, GET, PUT, DELETE)
// вызывает соответствующие обработчики
// если метод не поддерживается - возвращает статус 405 (Method Not Allowed)
func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost: // создание новой задачи
		addTaskHandler(w, r)
	case http.MethodGet: // получение информации о задаче
		getTaskHandler(w, r)
	case http.MethodPut: // редактирование существующей задачи
		editTaskHandler(w, r)
	case http.MethodDelete: // удаление задачи
		delTaskHandler(w, r)
	default:
		// если метод не поддерживается - возвращаем HTTP 405
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
