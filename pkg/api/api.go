package api

import (
	"net/http"
)

func Init() {
	http.HandleFunc("/api/nextdate", nextDateHandler)

	// auth-protected
	http.HandleFunc("/api/signin", signInHandler)

	http.HandleFunc("/api/tasks", auth(tasksHandler))
	http.HandleFunc("/api/task", auth(taskHandler))
	http.HandleFunc("/api/task/done", auth(doneHandler))
}
