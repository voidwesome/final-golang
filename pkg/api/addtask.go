package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"final-golang/pkg/db"
)

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var t db.Task
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		writeError(w, err, http.StatusBadRequest)
		return
	}
	// validations
	if t.Title == "" {
		writeError(w, errors.New("title required"), http.StatusBadRequest)
		return
	}
	if err := checkDate(&t); err != nil {
		writeError(w, err, http.StatusBadRequest)
		return
	}
	id, err := db.AddTask(&t)
	if err != nil {
		writeError(w, err, http.StatusInternalServerError)
		return
	}
	writeJSON(w, map[string]string{"id": fmt.Sprintf("%d", id)})
}

func checkDate(t *db.Task) error {
	now := time.Now()
	if t.Date == "" {
		t.Date = now.Format(DateLayout)
	}
	d, err := time.Parse(DateLayout, t.Date)
	if err != nil {
		return fmt.Errorf("bad date: %w", err)
	}
	if t.Repeat != "" {
		// validate and also get next date (we may need it)
		next, err := NextDate(now, t.Date, t.Repeat)
		if err != nil {
			return err
		}
		if !afterNow(d, now) {
			t.Date = next
		}
	} else {
		// no repeat
		if !afterNow(d, now) {
			t.Date = now.Format(DateLayout)
		}
	}
	return nil
}
