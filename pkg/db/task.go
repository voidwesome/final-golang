package db

import (
	"database/sql"
	"fmt"
)

// Task описывает структуру задачи
// ID - уникальный идентификатор задачи
// Date - дата выполнения задачи в формате YYYYMMDD
// Title - заголовок задачи
// Comment - дополнительный комментарий
// Repeat - правило повторения задачи
type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

// AddTask добавляет новую задачу в базу данных и возвращает её идентификатор
func AddTask(t *Task) (int64, error) {
	query := `INSERT INTO scheduler(date, title, comment, repeat) VALUES(?, ?, ?, ?)`
	res, err := DB.Exec(query, t.Date, t.Title, t.Comment, t.Repeat)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

// Tasks возвращает список задач, отсортированных по дате, ограниченный limit
func Tasks(limit int) ([]*Task, error) {
	query := `SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date ASC LIMIT ?`
	rows, err := DB.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close() // закрываем результат запроса после использования

	var list []*Task
	for rows.Next() {
		t := new(Task)
		// читаем данные из строки результата в структуру Task
		if err := rows.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat); err != nil {
			return nil, err
		}
		list = append(list, t)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	// если задач нет, возвращаем пустой список (а не nil)
	if list == nil {
		list = []*Task{}
	}
	return list, rows.Err()
}

// GetTask возвращает задачу по её идентификатору
func GetTask(id string) (*Task, error) {
	t := new(Task)
	err := DB.QueryRow(`SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?`, id).
		Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("задача не найдена")
	}
	return t, err
}

// UpdateTask обновляет все поля задачи (дата, заголовок, комментарий, правило повторения) по её ID
func UpdateTask(t *Task) error {
	query := `UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?`
	res, err := DB.Exec(query, t.Date, t.Title, t.Comment, t.Repeat, t.ID)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return fmt.Errorf("некорректный id для обновления задачи")
	}
	return nil
}

// UpdateDate обновляет только дату задачи по её ID
func UpdateDate(next string, id string) error {
	query := `UPDATE scheduler SET date = ? WHERE id = ?`
	res, err := DB.Exec(query, next, id)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return fmt.Errorf("некорректный id для обновления даты задачи")
	}
	return nil
}

// DeleteTask удаляет задачу по её ID
func DeleteTask(id string) error {
	res, err := DB.Exec(`DELETE FROM scheduler WHERE id = ?`, id)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return fmt.Errorf("задача не найдена")
	}
	return nil
}
