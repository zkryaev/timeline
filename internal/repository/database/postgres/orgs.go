package postgres

import (
	"context"
	"errors"
	"fmt"
	"timeline/internal/repository/models"

	"github.com/jackc/pgx/v5"
)

var (
	ErrOrgExists    = errors.New("org already exists")
	ErrOrgNotFound  = errors.New("org not found")
	ErrOrgsNotFound = errors.New("orgs not found")
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
		SELECT org_id, email, passwd_hash, name, rating, type, city, address, telephone, lat, long, about FROM orgs
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
		SELECT org_id, email, passwd_hash, name, rating, type, city, address, telephone, lat, long, about FROM orgs
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

/*
	Получение организаций
*/

// Принимает две точки на карте по диагонале. Возвращает список организаций лежащий в диапазоне этих двух точек
func (p *PostgresRepo) OrgsInArea(ctx context.Context, area *models.AreaParams) ([]*models.OrgSummary, error) {
	tx, err := p.db.Beginx()
	if err != nil {
		return nil, fmt.Errorf("failed to start tx: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	query := `SELECT org_id, name, rating, type FROM orgs
		WHERE lat BETWEEN $1 AND $2
		AND long BETWEEN $3 AND $4;
		`
	orgList := make([]*models.OrgSummary, 0, 1)
	if err = tx.SelectContext(
		ctx,
		&orgList,
		query,
		area.Left.Lat,
		area.Right.Lat,
		area.Left.Long,
		area.Right.Long,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrOrgsNotFound
		}
		return nil, fmt.Errorf("failed to get orgs at given area: %w", err)
	}
	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit tx: %w", err)
	}
	return orgList, nil
}

// Принимает пагинацию, имя и тип организации. Возвращает соответствующие организации
func (p *PostgresRepo) OrgsSearch(ctx context.Context, params *models.SearchParams) ([]*models.OrgInfo, error) {
	tx, err := p.db.Beginx()
	if err != nil {
		return nil, fmt.Errorf("failed to start tx: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	query := `SELECT org_id, name, rating, type, city, address, telephone, lat, long, about FROM orgs
		WHERE ($1 = '' OR ILIKE('%' || $1 || '%'))
		`
	orgList := make([]*models.OrgInfo, 0, params.Limit)
	if params.Type != "" {
		query += ` AND type = $2 LIMIT $3 OFFSET $4;`
		err = tx.SelectContext(ctx, &orgList, query, params.Name, params.Type, params.Limit, params.Offset)
	} else {
		query += ` LIMIT $2 OFFSET $3;`
		err = tx.SelectContext(ctx, &orgList, query, params.Name, params.Limit, params.Offset)
	}
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrOrgsNotFound
		}
		return nil, fmt.Errorf("failed orgs searching: %w", err)
	}
	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit tx: %w", err)
	}
	return orgList, nil
}
