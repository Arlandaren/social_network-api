package main

import (
	"log"
	"log/slog"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

func main() {
	logger := slog.Default()

	pgURL := os.Getenv("POSTGRES_CONN")
	if pgURL == "" {
		logger.Error("missed POSTGRES_CONN env")
		os.Exit(1)
	}

	db, err := sqlx.Connect("pgx", pgURL)
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	serverAddress := os.Getenv("SERVER_ADDRESS")
	if serverAddress == "" {
		logger.Error("missed SERVER_ADDRESS env (export smth like ':8080')")
		os.Exit(1)
	}

	s := NewServer(serverAddress, logger)

	err = s.Start()
	if err != nil {
		logger.Error("server has been stopped", "error", err)
	}
}
