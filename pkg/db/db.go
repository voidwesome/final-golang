package db

import (
	"database/sql"

	_ "modernc.org/sqlite"

	"os"
)

var DB *sql.DB

const schema = `
CREATE TABLE scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT, -- Уникальный идентификатор задачи
    date CHAR(8) NOT NULL DEFAULT "",     -- Дата задачи в формате YYYYMMDD
    title VARCHAR(256) NOT NULL DEFAULT "", -- Заголовок задачи
    comment TEXT NOT NULL DEFAULT "",     -- Комментарий к задаче
    repeat VARCHAR(128) NOT NULL DEFAULT "" -- Правило повторения задачи
);

CREATE INDEX idx_scheduler_date ON scheduler(date); -- Индекс для ускоренного поиска по дате
`

// Init инициализирует базу данных по указанному пути dbFile
// если база данных отсутствует, создается новая и применяется схема
func Init(dbFile string) error {
	install := false

	// проверяем, существует ли файл базы данных
	if _, err := os.Stat(dbFile); err != nil {
		install = true
	}

	// подключаемся к базе данных SQLite
	d, err := sql.Open("sqlite", dbFile)
	if err != nil {
		return err
	}
	DB = d
	// если база создается впервые, применяем схему
	if install {
		if _, err := DB.Exec(schema); err != nil {
			return err
		}
	}
	return nil
}
