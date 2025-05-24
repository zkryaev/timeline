//go:build integration

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
	cfg := &config.Database{
		Protocol: "postgres",
		Host:     "localhost",
		Port:     "5555",
		Name:     "testdb",
		User:     "user",
		Password: "passwd",
		SSLmode:  "disable",
	}
	db := New(cfg, "", true)
	if err := db.Open(); err != nil {
		log.Fatalf("SetupTest: \"%s\" cfg.presetDSN=\"%s\"", err.Error(), db.presetDSN)
	}
	suite.db = db
}

func TestPostgresTestSuite(t *testing.T) {
	suiter := &PostgresTestSuite{}
	suiter.SetupTest()
	defer suiter.db.Close()
	suite.Run(t, suiter)
}
