package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"final-golang/pkg/db"
)

func editTaskHandler(w http.ResponseWriter, r *http.Request) {
	var t db.Task
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		writeError(w, err, http.StatusBadRequest)
		return
	}
	if t.ID == "" {
		writeError(w, errors.New("id required"), http.StatusBadRequest)
		return
	}
	if t.Title == "" {
		writeError(w, errors.New("title required"), http.StatusBadRequest)
		return
	}
	if err := checkDate(&t); err != nil {
		writeError(w, err, http.StatusBadRequest)
		return
	}
	if err := db.UpdateTask(&t); err != nil {
		writeError(w, err, http.StatusNotFound)
		return
	}
	writeJSON(w, map[string]any{})
}
