package repository

import (
	"context"
	"fmt"
	"time"
	"timeline/internal/config"
	"timeline/internal/repository/database/postgres"
	"timeline/internal/repository/models"
)

type Database interface {
	Open() error
	Close()
	Repository
}

type Repository interface {
	CodeRepository
	UserRepository
	OrgRepository
}

type CodeRepository interface {
	SaveVerifyCode(ctx context.Context, info *models.CodeInfo) error
	VerifyCode(ctx context.Context, info *models.CodeInfo) (time.Time, error)
	ActivateAccount(ctx context.Context, id int, isOrg bool) error
	AccountExpiration(ctx context.Context, email string, isOrg bool) (*models.ExpInfo, error)
}

type UserRepository interface {
	UserSave(ctx context.Context, user *models.UserRegister) (int, error)
	UserByEmail(ctx context.Context, email string) (*models.UserRegister, error)
	UserByID(ctx context.Context, userID int) (*models.UserRegister, error)
}

type OrgRepository interface {
	OrgSave(ctx context.Context, org *models.OrgRegister) (int, error)
	OrgByEmail(ctx context.Context, email string) (*models.OrgRegister, error)
	OrgByID(ctx context.Context, id int) (*models.OrgRegister, error)
}

// Паттерн фабричный метод, чтобы не завязываться на конкретной БД
func GetDB(name string, cfg config.Database) (Database, error) {
	switch name {
	case "postgres":
		return postgres.New(cfg), nil
	default:
		return nil, fmt.Errorf("unexpected db name")
	}
}
