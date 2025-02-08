package main

import (
	"errors"
	"flag"
	"log"
	"time"

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

	var m *migrate.Migrate
	var err error
	// Две попытки между которыми 1 секунда
	maxRetries := 10
	for try := maxRetries; try > 0; try-- {
		m, err = migrate.New(
			"file://"+migrationsPath,
			dsn,
		)
		if err != nil {
			log.Printf("migrator: database is unavailable. Left %d attemps\n", try)
			time.Sleep(2 * time.Second)
		} else {
			break
		}
	}
	if err != nil {
		log.Fatal("create migrate: ", err.Error())
	}
	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Println("m.Up: no migrations to apply")
			return
		}
		log.Fatal("migrator: ", err)
	}
}
