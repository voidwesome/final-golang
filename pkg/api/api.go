package api

import (
	"log"
	"net/http"
)

func Init() {
	http.HandleFunc("/api/task", taskHandler)
	http.HandleFunc("/api/tasks", tasksHandler)
	http.HandleFunc("/api/task/done", MarkDoneHandler)
	http.HandleFunc("/api/nextdate", NextDateHandler)
}

func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		// Используем правильный обработчик
		postTaskHandler(w, r)
	case http.MethodGet:
		getTaskHandler(w, r)
	case http.MethodPut:
		UpdateTaskHandler(w, r)
	case http.MethodDelete:
		deleteTaskHandler(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
	log.Printf("Handled %v request", r)
}
