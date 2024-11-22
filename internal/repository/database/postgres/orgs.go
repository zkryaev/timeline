package postgres

import (
	"context"
	"errors"
	"fmt"
	"timeline/internal/repository/models"

	"github.com/jackc/pgx/v5"
)

var (
	ErrOrgExists   = errors.New("org already exists")
	ErrOrgNotFound = errors.New("org not found")
)

// Сохраняет в БД инфу об организации. Отдает org_id при успешном сохранении
// Вернет ошибку если такой пользователь уже существует
func (p *PostgresRepo) OrgSave(ctx context.Context, org *models.OrgRegister) (int, error) {
	tx, err := p.db.Beginx()
	if err != nil {
		return 0, fmt.Errorf("failed to start tx: %w", err)
	}
	// При возникновении ошибки транзакция откатывается по выходу из функции
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	query := `
		INSERT INTO orgs (email, passwd_hash, name, type, city, address, telephone, lat, long, about)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
        RETURNING org_id;
	`
	orgID := 0

	err = tx.QueryRowContext(ctx, query,
		org.HashCreds.Email,
		org.HashCreds.PasswdHash,
		org.Name,
		org.Type,
		org.City,
		org.Address,
		org.Telephone,
		org.Lat,
		org.Long,
		org.About,
	).Scan(&orgID)
	if err != nil {
		// TODO: пока отдаем фулл ошибку, а вообще нельзя внутренние ошибки отдавать наружу
		return 0, fmt.Errorf("failed to save org: %w", err)
	}
	if tx.Commit() != nil {
		return 0, fmt.Errorf("failed to commit transaction: %w", err)
	}
	return orgID, nil
}

func (p *PostgresRepo) OrgByEmail(ctx context.Context, email string) (*models.OrgRegister, error) {
	tx, err := p.db.Beginx()
	if err != nil {
		return nil, fmt.Errorf("failed to start tx: %w", err)
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	query := `
		SELECT org_id, email, passwd_hash, name, type, city, address, telephone, lat, long, about FROM orgs
        WHERE email = $1;
	`
	var org models.OrgRegister
	if err := tx.GetContext(ctx, &org, query, email); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrOrgNotFound
		}
		return nil, fmt.Errorf("failed get org by email: %w", err)
	}
	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit tx: %w", err)
	}

	return &org, nil
}

func (p *PostgresRepo) OrgByID(ctx context.Context, id int) (*models.OrgRegister, error) {
	tx, err := p.db.Beginx()
	if err != nil {
		return nil, fmt.Errorf("failed to start tx: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	query := `
		SELECT org_id, email, passwd_hash, name, type, city, address, telephone, lat, long, about FROM orgs
        WHERE org_id = $1;
	`
	var org models.OrgRegister
	if err = tx.GetContext(ctx, &org, query, id); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrOrgNotFound
		}
		return nil, fmt.Errorf("failed get org by id: %w", err)
	}
	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit tx: %w", err)
	}
	return &org, nil
}
