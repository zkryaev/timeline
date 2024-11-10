package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"timeline/internal/config"

	_ "github.com/lib/pq"
)

var (
	ErrUserExists   = errors.New("user already exists")
	ErrUserNotFound = errors.New("user not found")
)

type PostgresRepo struct {
	cfg config.Database
	db  *sql.DB
}

func New(cfg config.Database) *PostgresRepo {
	return &PostgresRepo{
		cfg: cfg,
	}
}

func (p *PostgresRepo) Open() error {
	db, err := sql.Open(p.cfg.Protocol, fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		p.cfg.Host,
		p.cfg.Port,
		p.cfg.User,
		p.cfg.Password,
		p.cfg.Name,
		p.cfg.SSLmode,
	))
	if err != nil {
		return err
	}
	if err = db.Ping(); err != nil {
		return err
	}
	p.db = db
	return nil
}

func (p *PostgresRepo) Close() {
	p.db.Close()
}
