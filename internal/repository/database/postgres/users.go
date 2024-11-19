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
	ErrUserExists   = errors.New("user already exists")
	ErrUserNotFound = errors.New("user not found")
	ErrCodeNotFound = errors.New("given code not found")
)

// Сохраняет пользователя и его почту/логин + пароль в хешированном виде
func (p *PostgresRepo) UserSave(ctx context.Context, user *models.UserRegisterModel) (int, error) {
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
		INSERT INTO users (email, passwd_hash, first_name, last_name, telephone, social, about)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
        RETURNING user_id;
	`
	userID := 0
	err = tx.QueryRowContext(
		ctx,
		query,
		user.HashCreds.Email,
		user.HashCreds.PasswdHash,
		user.FirstName,
		user.LastName,
		user.Telephone,
		user.Social,
		user.About,
	).Scan(&userID)
	if err != nil {
		return 0, fmt.Errorf("failed to save user: %w", err)
	}
	err = tx.Commit()
	if err != nil {
		return 0, fmt.Errorf("failed to commit tx: %w", err)
	}

	return userID, nil
}

// Получить юзера по email
func (p *PostgresRepo) UserByEmail(ctx context.Context, email string) (*entity.User, error) {
	tx, err := p.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to start tx: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	query := `
		SELECT user_id, email, passwd_hash, first_name, last_name, telephone, social, about FROM users
		WHERE email = $1;
	`
	var user entity.User
	var creds entity.HashCreds
	err = tx.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&creds.Email,
		&creds.PasswdHash,
		&user.Info.FirstName,
		&user.Info.LastName,
		&user.Info.Telephone,
		&user.Info.Social,
		&user.Info.About,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}
	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("failed to commit tx: %w", err)
	}
	return &user, nil
}

// Получить юзера по ID из БД
func (p *PostgresRepo) UserByID(ctx context.Context, id int) (*entity.User, error) {
	tx, err := p.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to start tx: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	query := `
		SELECT user_id, email, passwd_hash, user_name, telephone, social, about FROM users
		WHERE email = $1;
	`
	var user entity.User
	var creds entity.HashCreds
	err = tx.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&creds.Email,
		&creds.PasswdHash,
		&user.Info.FirstName,
		&user.Info.LastName,
		&user.Info.Telephone,
		&user.Info.Social,
		&user.Info.About,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to save user: %w", err)
	}
	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("failed to commit tx: %w", err)
	}
	return &user, nil
}

// Существуют ли юзер в БД
func (p *PostgresRepo) UserGetMetaInfo(ctx context.Context, email string) (*models.MetaInfo, error) {
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
        SELECT user_id, passwd_hash, created_at, verified FROM users
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
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to query select user: %w", err)
	}
	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("failed to commit tx: %w", err)
	}

	return resp, nil
}

func (p *PostgresRepo) UserIsExist(ctx context.Context, email string) (int, error) {
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
        SELECT user_id FROM users
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

// Сохранить код отправленный на почту пользователя
func (p *PostgresRepo) UserSaveCode(ctx context.Context, code string, user_id int) error {
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
		INSERT INTO user_verify (code, user_id)
        VALUES ($1, $2);
	`

	err = tx.QueryRowContext(
		ctx,
		query,
		code,
		user_id,
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

// Получить последний отправленный код на почту
// Если пришла ошибка значит веденный код неверен
func (p *PostgresRepo) UserCode(ctx context.Context, code string, user_id int) (time.Time, error) {
	tx, err := p.db.Begin()
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
		SELECT expires_at FROM user_verify
        WHERE code = $1 AND user_id = $2;
	`
	var expires_at time.Time
	err = tx.QueryRowContext(ctx, query, code, user_id).Scan(&expires_at)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return time.Time{}, ErrCodeNotFound
		}
		return time.Time{}, fmt.Errorf("failed to save user code: %w", err)
	}
	err = tx.Commit()
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to commit tx: %w", err)
	}
	return expires_at, nil
}

func (p *PostgresRepo) UserActivateAccount(ctx context.Context, user_id int) error {
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
		UPDATE users
		SET verified = $1
		WHERE user_id = $2;
	`
	var verified bool = true
	_, err = tx.ExecContext(ctx, query, verified, user_id)
	if err != nil {
		return fmt.Errorf("failed to activate user: %w", err)
	}
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit tx: %w", err)
	}

	return nil
}
