package postgres

import (
	"errors"
	"fmt"
	"timeline/internal/config"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

var (
	ErrNoRowsAffected = errors.New("no rows affected")
)

type PostgresRepo struct {
	cfg       *config.Database
	presetDSN string
	useCfg    bool
	db        *sqlx.DB
}

func New(cfg *config.Database, dsn string, useCfg bool) *PostgresRepo {
	return &PostgresRepo{
		cfg:       cfg,
		presetDSN: dsn,
		useCfg:    useCfg,
	}
}

func (p *PostgresRepo) Open() error {
	if p.useCfg {
		p.presetDSN = fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			p.cfg.Host,
			p.cfg.Port,
			p.cfg.User,
			p.cfg.Password,
			p.cfg.Name,
			p.cfg.SSLmode,
		)
	}
	db, err := sqlx.Connect("pgx", p.presetDSN)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	p.db = db
	return nil
}

func (p *PostgresRepo) Close() {
	p.db.Close()
}
