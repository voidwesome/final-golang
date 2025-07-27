package main

import (
	"final-golang/pkg/api"
	"final-golang/pkg/db"
	"log"
	"net/http"
	"os"
)

var port = os.Getenv("TODO_PORT")     // порт для запуска веб-сервера
var dbPath = os.Getenv("TODO_DBFILE") // путь к файлу базы данных

func main() {
	// если путь к БД не задан через переменную окружения, используем значение по умолчанию
	if dbPath == "" {
		dbPath = "scheduler.db"
	}

	// инициализируем подключение к базе данных
	if err := db.Init(dbPath); err != nil {
		log.Fatalf("Не удалось инициализировать БД: %v", err)
	}
	// закрываем соединение с базой данных при завершении программы
	defer db.DB.Close()

	api.DbConn = db.DB

	fs := http.FileServer(http.Dir("./web"))
	http.Handle("/", fs)

	// если порт не задан через переменную окружения, используем значение по умолчанию
	if port == "" {
		port = "7540"
	}
	addr := ":" + port

	// запускаем HTTP-сервер
	log.Printf("Сервер запущен на %s", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}
