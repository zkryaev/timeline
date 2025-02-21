package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
	"timeline/internal/infrastructure/models"

	"github.com/jackc/pgx/v5"
)

var (
	ErrCodeNotFound = errors.New("given code not found")
)

// Сохранить код отправленный на почту
func (p *PostgresRepo) SaveVerifyCode(ctx context.Context, Info *models.CodeInfo) error {
	tx, err := p.db.Beginx()
	if err != nil {
		return fmt.Errorf("failed to start tx: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	var query string
	switch Info.IsOrg {
	case false:
		query = `
		INSERT INTO users_verify (code, user_id)
        VALUES ($1, $2);
		`
	case true:
		query = `
		INSERT INTO orgs_verify (code, org_id)
        VALUES ($1, $2);
	`
	}

	if err = tx.QueryRowContext(ctx, query, Info.Code, Info.ID).Err(); err != nil {
		return fmt.Errorf("failed to save code: %w", err)
	}
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit tx: %w", err)
	}

	return nil
}

// Получает код и id, если находит код, возвращает его время сгорания
func (p *PostgresRepo) VerifyCode(ctx context.Context, Info *models.CodeInfo) (time.Time, error) {
	tx, err := p.db.Beginx()
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to start tx: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	var query string
	switch Info.IsOrg {
	case false:
		query = `
		SELECT expires_at FROM users_verify
        WHERE code = $1 AND user_id = $2;
		`
	case true:
		query = `
		SELECT expires_at FROM orgs_verify
        WHERE code = $1 AND org_id = $2;
	`
	}

	var expires_at time.Time
	err = tx.QueryRowContext(ctx, query, Info.Code, Info.ID).Scan(&expires_at)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return time.Time{}, ErrCodeNotFound
		}
		return time.Time{}, fmt.Errorf("failed to find code: %w", err)
	}
	if err = tx.Commit(); err != nil {
		return time.Time{}, fmt.Errorf("failed to commit tx: %w", err)
	}
	return expires_at, nil
}

// Устанавливает поле verified у сущности в true
func (p *PostgresRepo) ActivateAccount(ctx context.Context, id int, IsOrg bool) error {
	tx, err := p.db.Beginx()
	if err != nil {
		return fmt.Errorf("failed to start tx: %w", err)
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	var query string
	switch IsOrg {
	case false:
		query = `
		UPDATE users
		SET verified = $1
		WHERE is_delete = false 
		AND user_id = $2;
		`
	case true:
		query = `
		UPDATE orgs
		SET verified = $1
		WHERE is_delete = false 
		AND org_id = $2;
		`
	}
	res, err := tx.ExecContext(ctx, query, true, id)
	switch {
	case err != nil:
		return fmt.Errorf("failed to activate account: %w", err)
	default:
		if rowsAffected, _ := res.RowsAffected(); rowsAffected == 0{
			return ErrNoRowsAffected
		}
	}
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit tx: %w", err)
	}
	return nil
}

// Проверяет наличие организации в БД по ее почте/логину и возвращаем хеш пароля и ошибку
func (p *PostgresRepo) AccountExpiration(ctx context.Context, email string, IsOrg bool) (*models.ExpInfo, error) {
	tx, err := p.db.Beginx()
	if err != nil {
		return nil, fmt.Errorf("failed to start tx: %w", err)
	}
	// При возникновении ошибки транзакция откатывается по выходу из функции
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	var query string
	switch IsOrg {
	case false:
		query = `
        SELECT user_id, passwd_hash, created_at, verified FROM users
        WHERE is_delete = false 
		AND email = $1;
    	`
	case true:
		query = `
        SELECT org_id, passwd_hash, created_at, verified FROM orgs
        WHERE is_delete = false
		AND email = $1;
    	`
	}

	var Data models.ExpInfo
	if err = tx.QueryRowContext(ctx, query, email).Scan(&Data.ID, &Data.Hash, &Data.CreatedAt, &Data.Verified); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrOrgNotFound
		}
		return nil, fmt.Errorf("failed to get meta info: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit tx: %w", err)
	}

	return &Data, nil
}

// [CRON]: удаление стухших кодов
func (p *PostgresRepo) DeleteExpiredCodes(ctx context.Context) error {
	tx, err := p.db.Beginx()
	if err != nil {
		return fmt.Errorf("failed to start tx: %w", err)
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	query := `
		DELETE
		FROM users_verify
		WHERE expires_at <= (CURRENT_TIMESTAMP - INTERVAL '5 minute');

		DELETE
		FROM orgs_verify
		WHERE expires_at <= (CURRENT_TIMESTAMP - INTERVAL '5 minute');
		`
	res, err := tx.ExecContext(ctx, query)
	switch {
	case err != nil:
		return fmt.Errorf("failed to delete expired codes: %w", err)
	default:
		if rowsAffected, _ := res.RowsAffected(); rowsAffected == 0{
			return ErrNoRowsAffected
		}
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit tx: %w", err)
	}
	return nil
}
