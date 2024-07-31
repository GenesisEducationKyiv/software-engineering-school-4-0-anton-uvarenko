package main

import (
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	m, err := migrate.New(
		"file://./migrations/schema",
		"postgres://postgres:pass@127.0.0.1:5432/notification?sslmode=disable",
	)
	if err != nil {
		log.Fatal("can't run migrations: ", err)
	}

	err = m.Up()
	if err != nil {
		log.Fatal(err)
	}
}
