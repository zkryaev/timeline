package postgres_test

import (
	"log"
	"os"
	"testing"
	"timeline/internal/config"
	"timeline/internal/infrastructure"

	"github.com/stretchr/testify/suite"
)

type PostgresTestSuite struct {
	suite.Suite
	db infrastructure.Database
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
	db, err := infrastructure.GetDB(os.Getenv("DB"), cfg)
	if err != nil {
		log.Fatal(err.Error())
	}
	if err := db.Open(); err != nil {
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
