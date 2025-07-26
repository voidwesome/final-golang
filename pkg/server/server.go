package server

import (
	"fmt"
	"net/http"
	"os"

	"final-golang/pkg/api"
)

func Run() error {
	webDir := "web"
	if wd := os.Getenv("TODO_WEBDIR"); wd != "" {
		webDir = wd
	}

	api.Init()

	// file server for frontend
	http.Handle("/", http.FileServer(http.Dir(webDir)))

	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = "7540"
	}
	addr := ":" + port
	fmt.Printf("Listening on http://localhost%s\n", addr)
	return http.ListenAndServe(addr, nil)
}
