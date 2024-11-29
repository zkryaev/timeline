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
	UserUpdate(ctx context.Context, new *models.UserInfo) (*models.UserInfo, error)
	UserSave(ctx context.Context, user *models.UserRegister) (int, error)
	UserByID(ctx context.Context, userID int) (*models.UserInfo, error)
}

type OrgRepository interface {
	OrgSave(ctx context.Context, org *models.OrgRegister) (int, error)

	OrgByEmail(ctx context.Context, email string) (*models.OrgInfo, error)
	OrgByID(ctx context.Context, id int) (*models.OrgInfo, error)

	OrgsBySearch(ctx context.Context, params *models.SearchParams) ([]*models.OrgInfo, int, error)
	OrgsInArea(ctx context.Context, area *models.AreaParams) ([]*models.OrgSummary, error)

	OrgUpdate(ctx context.Context, new *models.OrgUpdate) (*models.OrgUpdate, error)
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
