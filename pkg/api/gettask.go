package api

import (
	"errors"
	"net/http"

	"github.com/voidwesome/final-golang/pkg/db"
)

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeError(w, errors.New("id required"), http.StatusBadRequest)
		return
	}
	t, err := db.GetTask(id)
	if err != nil {
		writeError(w, err, http.StatusNotFound)
		return
	}
	writeJSON(w, t)
}
