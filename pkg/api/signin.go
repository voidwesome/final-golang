package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
)

type signinReq struct {
	Password string `json:"password"`
}

func signInHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var req signinReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, err, http.StatusBadRequest)
		return
	}
	pass := os.Getenv("TODO_PASSWORD")
	if pass == "" {
		// аутентификация не требуется — но фронт её запрашивает
		// вернём фиктивный токен
		token, _ := makeToken("nopass")
		writeJSON(w, map[string]string{"token": token})
		return
	}
	if req.Password != pass {
		writeError(w, errors.New("Неверный пароль"), http.StatusUnauthorized)
		return
	}
	token, err := makeToken(pass)
	if err != nil {
		writeError(w, err, http.StatusInternalServerError)
		return
	}
	writeJSON(w, map[string]string{"token": token})
}
