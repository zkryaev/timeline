package main

import (
	"errors"
	"flag"
	"log"

	"github.com/golang-migrate/migrate/v4"
	// Driver for performing migrations into postgres
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	// Driver for getting migrations from files
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	var dsn, migrationsPath string
	flag.StringVar(&migrationsPath, "migrations-path", "", "path to the migrations folder")
	flag.StringVar(&dsn, "dsn", "", "database dsn (e.g., postgres://user:password@host:port/dbname?sslmode=disable)")
	flag.Parse()

	if dsn == "" {
		log.Fatal("dsn is required")
	}
	if migrationsPath == "" {
		log.Fatal("migrations-path is required")
	}

	m, err := migrate.New(
		"file://"+migrationsPath,
		dsn,
	)
	if err != nil {
		log.Fatal("migrate.New: ", err)
	}
	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Println("m.Up: no migrations to apply")
			return
		}
		log.Fatal("m.Up: ", err)
	}
}
