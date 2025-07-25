package main

import (
	"log"
	"net/http"
	"os"

	"final-golang/pkg/api"
	"final-golang/pkg/db"
)

var port = os.Getenv("TODO_PORT")
var dbPath = os.Getenv("TODO_DBFILE")

func main() {
	// Переменная окружения для БД
	if dbPath == "" {
		dbPath = "pkg/db/scheduler.db"
	}

	// Соединение с БД
	if err := db.Init(dbPath); err != nil {
		log.Fatalf("Не удалось инициализировать БД: %v", err)
	}

	// Передаем соединение в API
	api.DbConn = db.DB

	// Закрываем соединение с БД при завершении
	defer db.Close()

	// Регистрация эндпоинтов
	api.Init()

	// Настройка статики
	fs := http.FileServer(http.Dir("./web"))
	http.Handle("/", fs)

	// Запуск сервера
	if port == "" {
		port = "7540"
	}
	addr := ":" + port

	log.Printf("Сервер запущен на %s", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}
