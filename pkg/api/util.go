package api

import (
	"encoding/json"
	"net/http"
	"time"
)

const DateLayout = "20060102"

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	enc := json.NewEncoder(w)
	_ = enc.Encode(v)
}

func writeError(w http.ResponseWriter, err error, code int) {
	w.WriteHeader(code)
	writeJSON(w, map[string]string{"error": err.Error()})
}

// compare date-only (ignore time)
func afterNow(date, now time.Time) bool {
	d := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.Local)
	n := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
	return d.After(n)
}
