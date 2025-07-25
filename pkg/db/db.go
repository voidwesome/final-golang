package db

import (
	"database/sql"

	_ "modernc.org/sqlite"

	"os"
)

var DB *sql.DB

const schema = `
CREATE TABLE scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date CHAR(8) NOT NULL DEFAULT "",
    title VARCHAR(256) NOT NULL DEFAULT "",
    comment TEXT NOT NULL DEFAULT "",
    repeat VARCHAR(128) NOT NULL DEFAULT ""
);

CREATE INDEX idx_scheduler_date ON scheduler(date);
`

func Init(dbFile string) error {
	install := false
	if _, err := os.Stat(dbFile); err != nil {
		install = true
	}

	d, err := sql.Open("sqlite", dbFile)
	if err != nil {
		return err
	}
	DB = d

	if install {
		if _, err := DB.Exec(schema); err != nil {
			return err
		}
	}
	return nil
}
