package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
	"timeline/internal/entity"
	"timeline/internal/repository/models"
)

var (
	ErrOrgExists   = errors.New("org already exists")
	ErrOrgNotFound = errors.New("org not found")
)

// Сохраняет организацию + указанный города
func (p *PostgresRepo) OrgSave(ctx context.Context, org *models.OrgRegisterModel, cityName string) (int, error) {
	tx, err := p.db.BeginTx(ctx, &sql.TxOptions{})
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
	orgID, err := p.orgInfoSave(ctx, tx, org)
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

func (p *PostgresRepo) orgInfoSave(ctx context.Context, tx *sql.Tx, org *models.OrgRegisterModel) (int, error) {
	query := `
		INSERT INTO orgs (email, passwd_hash, org_name, org_address, telephone, social, about, lat, long, type)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
        RETURNING org_id;
	`
	orgID := 0

	err := tx.QueryRowContext(ctx, query,
		org.HashCreds.Email,
		org.HashCreds.PasswdHash,
		org.Name,
		org.Address,
		org.Telephone,
		org.Social,
		org.About,
		org.Lat,
		org.Long,
		org.Type,
	).Scan(&orgID)
	if err != nil {
		return 0, fmt.Errorf("failed to save org: %w", err)
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
func (p *PostgresRepo) OrgByEmail(ctx context.Context, email string) (*entity.Organization, error) {
	tx, err := p.db.BeginTx(ctx, &sql.TxOptions{ReadOnly: true})
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
		SELECT org_id, email, passwd_hash, org_name, org_address, telephone, social, about, lat, long FROM orgs
        WHERE email = $1;
	`
	var org entity.Organization
	var creds entity.HashCreds
	err = tx.QueryRowContext(ctx, query,
		email,
	).Scan(
		&org.ID,
		&creds.Email,
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
		return nil, fmt.Errorf("failed to save user: %w", err)
	}
	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("failed to commit tx: %w", err)
	}
	return &org, nil
}

// Получаем организацию по ее ID в БД
func (p *PostgresRepo) OrgByID(ctx context.Context, id int) (*entity.Organization, error) {
	tx, err := p.db.BeginTx(ctx, &sql.TxOptions{ReadOnly: true})
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
		SELECT org_id, email, passwd_hash, org_name, org_address, telephone, social, about, lat, long FROM orgs
        WHERE org_id = $1;
	`
	var org entity.Organization
	var creds entity.HashCreds
	err = tx.QueryRowContext(ctx, query, id).Scan(
		&org.ID,
		&creds.Email,
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
		return nil, fmt.Errorf("failed to save user: %w", err)
	}
	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("failed to commit tx: %w", err)
	}
	return &org, nil
}

// Проверяет наличие организации в БД по ее почте/логину и возвращаем хеш пароля и ошибку
func (p *PostgresRepo) OrgGetMetaInfo(ctx context.Context, email string) (*models.MetaInfo, error) {
	tx, err := p.db.BeginTx(ctx, &sql.TxOptions{ReadOnly: true})
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
        SELECT org_id, passwd_hash, created_at, verified FROM orgs
        WHERE email = $1;
    `

	resp := &models.MetaInfo{}
	err = tx.QueryRowContext(ctx, query, email).Scan(
		&resp.ID,
		&resp.Hash,
		&resp.CreatedAt,
		&resp.Verified,
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

	return resp, nil
}

func (p *PostgresRepo) OrgIsExist(ctx context.Context, email string) (int, error) {
	tx, err := p.db.BeginTx(ctx, &sql.TxOptions{})
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
        SELECT org_id FROM orgs
        WHERE email = $1;
    `
	var id int
	err = tx.QueryRowContext(ctx, query, email).Scan(
		&id,
	)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return 0, fmt.Errorf("failed to query select user: %w", err)
		}
	}
	errtx := tx.Commit()
	if errtx != nil {
		return 0, fmt.Errorf("failed to commit tx: %w", err)
	}

	return id, nil
}

// Сохраняем код отправленный на почту организации
func (p *PostgresRepo) OrgSaveCode(ctx context.Context, code string, org_id int) error {
	tx, err := p.db.BeginTx(ctx, &sql.TxOptions{})
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
		INSERT INTO orgs_verify (code, org_id, expires_at)
        VALUES ($1, $2, $3);
	`

	err = tx.QueryRowContext(
		ctx,
		query,
		code,
		org_id,
		time.Now(),
	).Err()
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
func (p *PostgresRepo) OrgCode(ctx context.Context, code string, org_id int) (time.Time, error) {
	tx, err := p.db.BeginTx(ctx, &sql.TxOptions{ReadOnly: true})
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to start tx: %w", err)
	}
	// При возникновении ошибки транзакция откатывается по выходу из функции
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	query := `
		SELECT expires_at FROM orgs_verify
        WHERE code = $1 AND org_id = $2;
	`
	var expires_at time.Time
	err = tx.QueryRowContext(ctx, query, code, org_id).Scan(&expires_at)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to save org code: %w", err)
	}
	err = tx.Commit()
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to commit tx: %w", err)
	}
	expires_at.Add(3 * time.Hour) // т.к время по гринвичу отстает на 3 часа от МСК
	return expires_at, nil
}

func (p *PostgresRepo) OrgActivateAccount(ctx context.Context, org_id int) error {
	tx, err := p.db.BeginTx(ctx, &sql.TxOptions{})
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
		UPDATE orgs
		SET verified = $1
		WHERE org_id = $2;
	`
	var verified bool = true
	_, err = tx.ExecContext(ctx, query, verified, org_id)
	if err != nil {
		return fmt.Errorf("failed to activate org: %w", err)
	}
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit tx: %w", err)
	}

	return nil
}
