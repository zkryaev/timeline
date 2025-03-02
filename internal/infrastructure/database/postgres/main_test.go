package postgres

import (
	"log"
	"testing"
	"timeline/internal/config"

	"github.com/stretchr/testify/suite"
)

type PostgresTestSuite struct {
	suite.Suite
	db *PostgresRepo
}

func (suite *PostgresTestSuite) SetupTest() {
	cfg := config.Database{
		Protocol: "postgres",
		Host:     "localhost",
		Port:     "5555",
		Name:     "testdb",
		User:     "user",
		Password: "passwd",
		SSLmode:  "disable",
	}
	db := New(cfg)
	err := db.Open()
	if err != nil {
		log.Fatal(err.Error())
	}
	if err = db.Open(); err != nil {
		log.Fatal(err.Error())
	}
	suite.db = db
}

func TestPostgresTestSuite(t *testing.T) {
	suiter := &PostgresTestSuite{}
	suiter.SetupTest()
	defer suiter.db.Close()
	suite.Run(t, suiter)
}
