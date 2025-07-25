package main

import (
	"log"
	"os"

	"final-golang/pkg/db"
	"final-golang/pkg/server"
)

func main() {
	dbFile := os.Getenv("TODO_DBFILE")
	if dbFile == "" {
		dbFile = "scheduler.db"
	}
	if err := db.Init(dbFile); err != nil {
		log.Fatalf("db init error: %v", err)
	}
	if err := server.Run(); err != nil {
		log.Fatalf("server run error: %v", err)
	}
}
