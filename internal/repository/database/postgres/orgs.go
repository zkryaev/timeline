package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"timeline/internal/entity"
	"timeline/internal/repository/models"
)

var (
	ErrOrgExists   = errors.New("org already exists")
	ErrOrgNotFound = errors.New("org not found")
)

// Сохраняет организацию + указанный города
func (p *PostgresRepo) OrgSave(ctx context.Context, org *entity.OrgInfo, creds *entity.Credentials, cityName string) (int, error) {
	tx, err := p.db.Begin()
	if err != nil {
		return 0, fmt.Errorf("failed to start tx: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Сначала сохраняем город и получаем его ID
	cityID, err := orgSaveCity(ctx, tx, cityName)
	if err != nil {
		return 0, fmt.Errorf("failed to save city: %w", err)
	}

	// Сохраняем организацию
	orgID, err := p.orgInfoSave(ctx, org, creds)
	if err != nil {
		return 0, fmt.Errorf("failed to save org: %w", err)
	}

	// Связываем организацию с городом
	err = orgCityLink(ctx, tx, orgID, cityID)
	if err != nil {
		return 0, fmt.Errorf("failed to link org to city: %w", err)
	}

	// Коммитим транзакцию
	err = tx.Commit()
	if err != nil {
		return 0, fmt.Errorf("failed to commit tx: %w", err)
	}

	return orgID, nil
}

func (p *PostgresRepo) orgInfoSave(ctx context.Context, org *entity.OrgInfo, creds *entity.Credentials) (int, error) {
	tx, err := p.db.Begin()
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
		INSERT INTO orgs (email, passwd_hash, org_name, org_address, telephone, social, about, lat, long)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
        RETURNING org_id;
	`
	orgID := 0
	err = tx.QueryRowContext(ctx, query,
		creds.Login,
		creds.PasswdHash,
		org.Name,
		org.Address,
		org.Telephone,
		org.Social,
		org.About,
		org.Lat,
		org.Long,
	).Scan(&orgID)
	if err != nil {
		return 0, fmt.Errorf("failed to save org: %w", err)
	}
	err = tx.Commit()
	if err != nil {
		return 0, fmt.Errorf("failed to commit tx: %w", err)
	}

	return orgID, nil
}

func orgSaveCity(ctx context.Context, tx *sql.Tx, cityName string) (int, error) {
	var cityID int
	query := `
		INSERT INTO city (name)
		VALUES ($1)
		RETURNING city_id;
	`
	err := tx.QueryRowContext(ctx, query, cityName).Scan(&cityID)
	if err != nil {
		return 0, fmt.Errorf("failed to save city: %w", err)
	}
	return cityID, nil
}

func orgCityLink(ctx context.Context, tx *sql.Tx, orgID, cityID int) error {
	query := `
		INSERT INTO orgs_city (org_id, city_id)
		VALUES ($1, $2);
	`
	_, err := tx.ExecContext(ctx, query, orgID, cityID)
	if err != nil {
		return fmt.Errorf("failed to link org to city: %w", err)
	}
	return nil
}

// Получить организацию по ее почте/логину
func (p *PostgresRepo) OrgByEmail(ctx context.Context, email string) (*entity.Organization, *entity.Credentials, error) {
	tx, err := p.db.Begin()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to start tx: %w", err)
	}
	// При возникновении ошибки транзакция откатывается по выходу из функции
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	query := `
		SELECT email, passwd_hash, org_name, org_address, telephone, social, about, lat, long FROM orgs
        WHERE email = $1;
	`
	var org entity.Organization
	var creds entity.Credentials
	err = tx.QueryRowContext(ctx, query,
		email,
	).Scan(
		&creds.Login,
		&creds.PasswdHash,
		&org.Info.Name,
		&org.Info.Address,
		&org.Info.Telephone,
		&org.Info.Social,
		&org.Info.About,
		&org.Info.Lat,
		&org.Info.Long,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to save user: %w", err)
	}
	err = tx.Commit()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to commit tx: %w", err)
	}
	return &org, &creds, nil
}

// Получаем организацию по ее ID в БД
func (p *PostgresRepo) OrgByID(ctx context.Context, id int) (*entity.Organization, *entity.Credentials, error) {
	tx, err := p.db.Begin()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to start tx: %w", err)
	}
	// При возникновении ошибки транзакция откатывается по выходу из функции
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	query := `
		SELECT email, passwd_hash, org_name, org_address, telephone, social, about, lat, long FROM orgs
        WHERE org_id = $1;
	`
	var org entity.Organization
	var creds entity.Credentials
	err = tx.QueryRowContext(ctx, query,
		id,
	).Scan(
		&creds.Login,
		&creds.PasswdHash,
		&org.Info.Name,
		&org.Info.Address,
		&org.Info.Telephone,
		&org.Info.Social,
		&org.Info.About,
		&org.Info.Lat,
		&org.Info.Long,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to save user: %w", err)
	}
	err = tx.Commit()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to commit tx: %w", err)
	}
	return &org, &creds, nil
}

// Проверяет наличие организации в БД по ее почте/логину и возвращаем хеш пароля и ошибку
func (p *PostgresRepo) OrgIsExist(ctx context.Context, email string) (*models.IsExistResponse, error) {
	tx, err := p.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to start tx: %w", err)
	}
	// При возникновении ошибки транзакция откатывается по выходу из функции
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	query := `
        SELECT org_id, passwd_hash, created_at FROM orgs
        WHERE email = $1;
    `

	var MetaInfo models.IsExistResponse
	err = tx.QueryRowContext(ctx, query, email).Scan(
		&MetaInfo.ID,
		&MetaInfo.Hash,
		&MetaInfo.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrOrgNotFound
		}
		return nil, fmt.Errorf("failed to query select org: %w", err)
	}
	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("failed to commit tx: %w", err)
	}

	return &MetaInfo, nil
}

// Сохраняем код отправленный на почту организации
func (p *PostgresRepo) OrgSaveCode(ctx context.Context, code string, org_id int) error {
	tx, err := p.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to start tx: %w", err)
	}
	// При возникновении ошибки транзакция откатывается по выходу из функции
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	query := `
		INSERT INTO org_verify (code, org_id)
        VALUES ($1, $2);
	`

	err = tx.QueryRowContext(ctx, query, code, org_id).Err()
	if err != nil {
		return fmt.Errorf("failed to save org code: %w", err)
	}
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit tx: %w", err)
	}

	return nil
}

// Получаем последний код отправленный на почту организации
// Если ошибка, значит веденный код неверный
func (p *PostgresRepo) OrgCode(ctx context.Context, code string, org_id int) (string, error) {
	tx, err := p.db.Begin()
	if err != nil {
		return "", fmt.Errorf("failed to start tx: %w", err)
	}
	// При возникновении ошибки транзакция откатывается по выходу из функции
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	query := `
		SELECT code FROM org_verify
        WHERE code = $1 AND org_id = $2;
	`
	var verifyCode string
	err = tx.QueryRowContext(ctx, query, code, org_id).Scan(&verifyCode)
	if err != nil {
		return "", fmt.Errorf("failed to save org code: %w", err)
	}
	err = tx.Commit()
	if err != nil {
		return "", fmt.Errorf("failed to commit tx: %w", err)
	}

	return verifyCode, nil
}
