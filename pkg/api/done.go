package api

import (
	"net/http"
	"time"

	"final-golang/pkg/db"
)

// MarkDoneHandler обрабатывает POST-запросы для отметки задачи как выполненной
// если у задачи есть правило повторения, обновляет дату на следующую по правилу
// иначе - удаляет задачу из базы.
func MarkDoneHandler(w http.ResponseWriter, r *http.Request) {
	// Проверяем, что метод запроса POST
	if r.Method != http.MethodPost {
		writeJSON(w, errorResponse{Error: "Метод не разрешён"})
		return
	}

	// получаем ID задачи из параметров запроса
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSON(w, errorResponse{Error: "Не указан ID"})
		return
	}

	// получаем задачу из базы данных
	t, err := db.GetTask(id)
	if err != nil {
		writeJSON(w, errorResponse{Error: "Задача не найдена"})
		return
	}

	// если правило повторения отсутствует - удаляем задачу
	if t.Repeat == "" {
		if err := db.DeleteTask(id); err != nil {
			writeJSON(w, errorResponse{Error: err.Error()})
			return
		}
		writeJSON(w, struct{}{}) // Пустой успешный ответ
		return
	}

	// парсим текущую дату задачи для вычисления следующей даты
	prev, err := time.Parse(dateLayout, t.Date)
	if err != nil {
		writeJSON(w, errorResponse{Error: "Некорректная дата в задаче"})
		return
	}

	// вычисляем следующую дату выполнения задачи на основе правила повторения.
	nextDate, err := NextDate(prev, t.Date, t.Repeat)
	if err != nil {
		writeJSON(w, errorResponse{Error: err.Error()})
		return
	}

	// обновляем дату задачи в базе данных
	if err := db.UpdateDate(nextDate, id); err != nil {
		writeJSON(w, errorResponse{Error: err.Error()})
		return
	}

	writeJSON(w, struct{}{}) // пустой успешный ответ
}
