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
	var dsn, migrationsPath, instanceName string
	flag.StringVar(&migrationsPath, "migrations-path", "", "path to the migrations folder")
	flag.StringVar(&dsn, "dsn", "", "database dsn (e.g., postgres://user:password@host:port/dbname?sslmode=disable)")
	flag.StringVar(&instanceName, "instance", "", "instance name (e.g., analytics_db)")
	flag.Parse()

	if dsn == "" {
		log.Fatalf("instance: %s: dsn is required", instanceName)
	}
	if migrationsPath == "" {
		log.Fatalf("instance: %s: migrations-path is required", instanceName)
	}

	var m *migrate.Migrate
	var err error

	maxRetries := 10
	for try := maxRetries; try > 0; try-- {
		m, err = migrate.New(
			"file://"+migrationsPath,
			dsn,
		)
		if err != nil {
			log.Printf("migrator: instance: %s: failed open database. Left %d attemps\n", instanceName, try)
			time.Sleep(2 * time.Second)
		} else {
			break
		}
	}
	if err != nil {
		log.Fatalf("migrator: instance: %s: %s", instanceName, err.Error())
	}
	if err = m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Printf("migrator: instance: %s: no migrations to apply", instanceName)
			return
		}
		log.Fatalf("migrator: instance: %s: %s", instanceName, err.Error())
	}
}
