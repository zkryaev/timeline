package postgres

import (
	"context"
	"fmt"
	"timeline/internal/config"

	_ "github.com/jackc/pgx"
	"github.com/jmoiron/sqlx"
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
	context.Background()
	db, err := sqlx.Connect(p.cfg.Protocol, fmt.Sprintf(
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
