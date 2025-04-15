package main

import (
	"flag"
	"fmt"
	"log"
	"timeline/internal/infrastructure/database/postgres"
	"timeline/internal/utils/loader"

	"go.uber.org/zap"
)

func main() {
	var givenDSN, instanceName string
	flag.StringVar(&givenDSN, "dsn", "", "database dsn (e.g., postgres://user:password@host:port/dbname?sslmode=disable)")
	flag.StringVar(&instanceName, "instance", "not set", "instance name (e.g., test_db)")
	flag.Parse()
	logger, err := zap.NewDevelopmentConfig().Build()
	if err != nil {
		log.Fatal("loader: ", err.Error())
	}
	if givenDSN == "" {
		logger.Fatal(fmt.Sprintf("instance: %s: dsn is required", instanceName))
	}
	db := postgres.New(nil, givenDSN, false)
	if err := db.Open(); err != nil {
		logger.Fatal(fmt.Sprintf("instance: %s: failed to connect to db", instanceName))
	}
	defer db.Close()
	backdata := &loader.BackData{}
	loader.LoadData(logger, db, backdata)
}
