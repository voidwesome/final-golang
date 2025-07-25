package db

import (
	"database/sql"
	"errors"
	"fmt"
	"os"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

const schema = `
CREATE TABLE IF NOT EXISTS scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date CHAR(8) NOT NULL DEFAULT "",
    title VARCHAR(255) NOT NULL DEFAULT "",
    comment TEXT NOT NULL DEFAULT "",
    repeat VARCHAR(128) NOT NULL DEFAULT ""
);

CREATE INDEX IF NOT EXISTS idx_scheduler_date ON scheduler(date);
`

func Close() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}

func Init(dbFile string) error {
	install := false
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		install = true
	}

	var err error
	DB, err = sql.Open("sqlite", dbFile)
	if err != nil {
		return fmt.Errorf("не удалось открыть БД: %w", err)
	}

	if install {
		_, err := DB.Exec(schema)
		if err != nil {
			return fmt.Errorf("ошибка создания таблицы: %w", err)
		}
	}

	return nil
}

func Tasks(limit int) ([]*Task, error) {
	rows, err := DB.Query("SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT ?", limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*Task
	for rows.Next() {
		var t Task
		err := rows.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, &t)
	}
	return tasks, nil
}

func GetTask(id string) (*Task, error) {
	var t Task
	err := DB.QueryRow("SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?", id).
		Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func DeleteTask(id string) error {
	_, err := DB.Exec("DELETE FROM scheduler WHERE id = ?", id)
	return err
}

func UpdateDate(date string, id string) error {
	res, err := DB.Exec("UPDATE scheduler SET date = ? WHERE id = ?", date, id)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return errors.New("задача не найдена")
	}
	return nil
}
