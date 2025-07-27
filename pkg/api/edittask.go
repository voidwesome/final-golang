package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"final-golang/pkg/db"
)

// editTaskHandler обрабатывает HTTP-запрос на редактирование задачи
// ожидает JSON с данными задачи в теле запроса
func editTaskHandler(w http.ResponseWriter, r *http.Request) {
	var t db.Task

	// декодируем JSON из тела запроса в структуру задачи
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		writeError(w, err, http.StatusBadRequest) // ошибка десериализации
		return
	}

	// проверяем, что указан ID задачи
	if t.ID == "" {
		writeError(w, errors.New("требуется id"), http.StatusBadRequest) // ошибка отсутствия ID
		return
	}

	// проверяем, что указан заголовок задачи
	if t.Title == "" {
		writeError(w, errors.New("требуется заголовок"), http.StatusBadRequest) // ошибка отсутствия заголовка
		return
	}

	// проверяем корректность даты задачи
	if err := checkDate(&t); err != nil {
		writeError(w, err, http.StatusBadRequest) //	 ошибка некорректной даты
		return
	}

	// обновляем задачу в базе данных
	if err := db.UpdateTask(&t); err != nil {
		writeError(w, err, http.StatusNotFound) // ошибка обновления задачи
		return
	}

	// отправляем пустой JSON-ответ об успешном обновлении
	writeJSON(w, map[string]any{})
}
