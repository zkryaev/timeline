package postgres

import (
	"context"
	"database/sql"
	"errors"
	"timeline/internal/model"
)

var (
	ErrUserExists   = errors.New("user already exists")
	ErrUserNotFound = errors.New("user not found")
)

type PostgresRepo struct {
	DB *sql.DB // Указатель на соединение с базой данных
}

// TODO: что передавать в БД?
func (p *PostgresRepo) SaveUser(ctx context.Context, user *model.User) (uint64, error) { return 0, nil }

func (p *PostgresRepo) User(ctx context.Context) (*model.User, error) { return nil, nil }

func (p *PostgresRepo) SaveOrg(ctx context.Context, org *model.Organization) (uint64, error) {
	return 0, nil
}

func (p *PostgresRepo) Organization(ctx context.Context) (*model.Organization, error) {
	return nil, nil
}
