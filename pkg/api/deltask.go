package api

import (
	"errors"
	"net/http"

	"final-golang/pkg/db"
)

func delTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeError(w, errors.New("id required"), http.StatusBadRequest)
		return
	}
	if err := db.DeleteTask(id); err != nil {
		writeError(w, err, http.StatusNotFound)
		return
	}
	writeJSON(w, map[string]any{})
}
