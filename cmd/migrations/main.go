package main

import (
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	m, err := migrate.New(
		"file:///home/migrations/schema",
		"postgres://postgres:pass@db:5432/postgres?sslmode=disable",
	)
	if err != nil {
		log.Fatal("can't run migrations: ", err)
	}

	err = m.Up()
	if err != nil {
		log.Fatal(err)
	}
}
