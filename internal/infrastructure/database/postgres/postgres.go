package postgres

import (
	"errors"
	"fmt"
	"timeline/internal/config"

	"github.com/jmoiron/sqlx"
)

var (
	ErrNoRowsAffected = errors.New("no rows affected")
)

type PostgresRepo struct {
	cfg config.Database
	db  *sqlx.DB
}

func New(cfg config.Database) *PostgresRepo {
	return &PostgresRepo{
		cfg: cfg,
	}
}

func (p *PostgresRepo) Open() error {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		p.cfg.Host,
		p.cfg.Port,
		p.cfg.User,
		p.cfg.Password,
		p.cfg.Name,
		p.cfg.SSLmode,
	)
	db, err := sqlx.Connect("pgx", dsn)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
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
