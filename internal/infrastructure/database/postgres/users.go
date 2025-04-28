package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"timeline/internal/infrastructure/models/usermodel"
)

var (
	ErrUserExists   = errors.New("user already exists")
	ErrUserNotFound = errors.New("user not found")
)

// Принимает всю регистрационную инфу. Возвращает user_id
// Если такой пользователь уже существует -> ошибка
func (p *PostgresRepo) UserSave(ctx context.Context, user *usermodel.UserRegister) (int, error) {
	tx, err := p.db.Beginx()
	if err != nil {
		return 0, fmt.Errorf("failed to start tx: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	query := `
		INSERT INTO users (email, passwd_hash, first_name, last_name, telephone, city, about)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
        RETURNING user_id;
	`
	var userID int
	if err = tx.QueryRowContext(
		ctx,
		query,
		user.HashCreds.Email,
		user.HashCreds.PasswdHash,
		user.FirstName,
		user.LastName,
		user.Telephone,
		user.City,
		user.About,
	).Scan(&userID); err != nil {
		return 0, fmt.Errorf("failed to save user: %w", err)
	}
	if err = tx.Commit(); err != nil {
		return 0, fmt.Errorf("failed to commit tx: %w", err)
	}
	return userID, nil
}

func (p *PostgresRepo) UserByID(ctx context.Context, userID int) (*usermodel.UserInfo, error) {
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
		SELECT user_id, email, COALESCE(uuid, '') AS uuid, first_name, last_name, city, telephone, about FROM users
		WHERE is_delete = false 
		AND user_id = $1;
	`
	var user usermodel.UserInfo
	if err = tx.GetContext(ctx, &user, query, userID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}
	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit tx: %w", err)
	}
	return &user, nil
}

func (p *PostgresRepo) UserUpdate(ctx context.Context, user *usermodel.UserInfo) error {
	tx, err := p.db.Beginx()
	if err != nil {
		return fmt.Errorf("failed to start tx: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	query := `UPDATE users
		SET 
			first_name = $1,
			last_name = $2,
			city = $3,
			telephone = $4,
			about = $5
		WHERE is_delete = false
		AND user_id = $6;
		`
	res, err := tx.ExecContext(ctx, query,
		user.FirstName,
		user.LastName,
		user.City,
		user.Telephone,
		user.About,
		user.UserID,
	)
	switch {
	case err != nil:
		return fmt.Errorf("failed update user's info: %w", err)
	default:
		if rowsAffected, _ := res.RowsAffected(); rowsAffected == 0 {
			return ErrNoRowsAffected
		}
	}
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit tx: %w", err)
	}
	return nil
}

// [CRON]: удаление стухших юзер аккаунтов
func (p *PostgresRepo) UserDeleteExpired(ctx context.Context) error {
	tx, err := p.db.Beginx()
	if err != nil {
		return fmt.Errorf("failed to start tx: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	query := `DELETE FROM users
		WHERE verified IS NOT true
		AND created_at < CURRENT_TIMESTAMP - INTERVAL '1 day';
	`
	res, err := tx.ExecContext(ctx, query)
	switch {
	case err != nil:
		return fmt.Errorf("failed to delete expired user's accounts: %w", err)
	default:
		if rowsAffected, _ := res.RowsAffected(); rowsAffected == 0 {
			return ErrNoRowsAffected
		}
	}
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit tx: %w", err)
	}
	return nil
}

func (p *PostgresRepo) UserSoftDelete(ctx context.Context, userID int) error {
	tx, err := p.db.Beginx()
	if err != nil {
		return fmt.Errorf("failed to start tx: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	query := `UPDATE users
		SET 
			is_delete = TRUE
		WHERE user_id = $1;
	`
	res, err := tx.ExecContext(ctx, query, userID)
	switch {
	case err != nil:
		return fmt.Errorf("failed to delete user: %w", err)
	default:
		if rowsAffected, _ := res.RowsAffected(); rowsAffected == 0 {
			return ErrNoRowsAffected
		}
	}
	query = `
		DELETE FROM users_verify 
		WHERE user_id = $1;
	`
	_, err = tx.ExecContext(ctx, query, userID)
	switch {
	case err != nil:
		return fmt.Errorf("failed to delete user: %w", err)
	default:
		if rowsAffected, _ := res.RowsAffected(); rowsAffected == 0 {
			return ErrNoRowsAffected
		}
	}
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit tx: %w", err)
	}
	return nil
}

// Only for tests!
func (p *PostgresRepo) UserDelete(ctx context.Context, userID int) error {
	tx, err := p.db.Beginx()
	if err != nil {
		return fmt.Errorf("failed to start tx: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	query := `DELETE FROM users
		WHERE user_id = $1;
	`
	res, err := tx.ExecContext(ctx, query, userID)
	switch {
	case err != nil:
		return fmt.Errorf("failed to delete user: %w", err)
	default:
		if rowsAffected, _ := res.RowsAffected(); rowsAffected == 0 {
			return ErrNoRowsAffected
		}
	}
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit tx: %w", err)
	}
	return nil
}

func (p *PostgresRepo) UserUUID(ctx context.Context, userID int) (string, error) {
	tx, err := p.db.Beginx()
	if err != nil {
		return "", fmt.Errorf("failed to start tx: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	query := `
		SELECT COALESCE(uuid, '') AS uuid
		FROM users
		WHERE user_id = $1;
	`
	var uuid string
	if err = tx.QueryRowContext(ctx, query, userID).Scan(&uuid); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", ErrUserNotFound
		}
		return "", fmt.Errorf("failed to get user uuid by id: %w", err)
	}
	if err = tx.Commit(); err != nil {
		return "", fmt.Errorf("failed to commit tx: %w", err)
	}
	return uuid, nil
}

func (p *PostgresRepo) UserSetUUID(ctx context.Context, userID int, newUUID string) error {
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
		UPDATE users
		SET
			uuid = $1
		WHERE user_id = $2;
	`
	res, err := tx.ExecContext(ctx, query, newUUID, userID)
	switch {
	case err != nil:
		return fmt.Errorf("failed to set user's uuid: %w", err)
	default:
		if rowsAffected, _ := res.RowsAffected(); rowsAffected == 0 {
			return ErrNoRowsAffected
		}
	}
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit tx: %w", err)
	}
	return nil
}

func (p *PostgresRepo) UserDeleteURL(ctx context.Context, userID int, url string) error {
	tx, err := p.db.Beginx()
	if err != nil {
		return fmt.Errorf("failed to start tx: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	query := `UPDATE users SET uuid = '' WHERE uuid = $1 AND user_id = $2;`
	res, err := tx.ExecContext(ctx, query, url, userID)
	switch {
	case err != nil:
		return fmt.Errorf("failed to delete user's url: %w", err)
	default:
		if rowsAffected, _ := res.RowsAffected(); rowsAffected == 0 {
			return ErrNoRowsAffected
		}
	}
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit tx: %w", err)
	}
	return nil
}
