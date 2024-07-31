package db

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
)

func Connect() *pgx.Conn {
	db, err := pgx.Connect(context.Background(), os.Getenv("CONNECTION_STRING"))
	if err != nil {
		log.Fatal(err)
	}

	return db
}
