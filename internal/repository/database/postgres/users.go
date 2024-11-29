package postgres

import (
	"context"
	"errors"
	"fmt"
	"timeline/internal/repository/models"

	"github.com/jackc/pgx/v5"
)

var (
	ErrUserExists   = errors.New("user already exists")
	ErrUserNotFound = errors.New("user not found")
)

// Принимает всю регистрационную инфу. Возвращает user_id
// Если такой пользователь уже существует -> ошибка
func (p *PostgresRepo) UserSave(ctx context.Context, user *models.UserRegister) (int, error) {
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
		INSERT INTO users (email, passwd_hash, first_name, last_name, telephone, city, about)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
        RETURNING user_id;
	`
	var UserID int
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
	).Scan(&UserID); err != nil {
		return 0, fmt.Errorf("failed to save user: %w", err)
	}
	if err = tx.Commit(); err != nil {
		return 0, fmt.Errorf("failed to commit tx: %w", err)
	}

	return UserID, nil
}

func (p *PostgresRepo) UserByID(ctx context.Context, UserID int) (*models.UserInfo, error) {
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
		SELECT user_id, first_name, last_name, city, telephone, about FROM users
		WHERE user_id = $1;
	`
	var user models.UserInfo
	if err = tx.GetContext(ctx, &user, query, UserID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}
	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit tx: %w", err)
	}
	return &user, nil
}

func (p *PostgresRepo) UserUpdate(ctx context.Context, new *models.UserInfo) error {
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
		WHERE user_id = $6
		`
	res, err := tx.ExecContext(ctx, query,
		new.FirstName,
		new.LastName,
		new.City,
		new.Telephone,
		new.About,
		new.UserID,
	)
	if err != nil {
		return fmt.Errorf("failed update user info: %w", err)
	}
	if res != nil {
		if _, errNoRowsAffected := res.RowsAffected(); errNoRowsAffected != nil {
			return ErrOrgNotFound
		}
	}
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit tx: %w", err)
	}
	return nil
}
