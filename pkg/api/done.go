package api

import (
	"errors"
	"net/http"
	"time"

	"final-golang/pkg/db"
)

func doneHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
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
	// if no repeat -> delete
	if t.Repeat == "" {
		if err := db.DeleteTask(id); err != nil {
			writeError(w, err, http.StatusInternalServerError)
			return
		}
		writeJSON(w, map[string]any{})
		return
	}
	now := time.Now()
	next, err := NextDate(now, t.Date, t.Repeat)
	if err != nil {
		writeError(w, err, http.StatusBadRequest)
		return
	}
	if err := db.UpdateDate(next, id); err != nil {
		writeError(w, err, http.StatusInternalServerError)
		return
	}
	writeJSON(w, map[string]any{})
}
