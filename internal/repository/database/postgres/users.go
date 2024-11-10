package postgres

import (
	"context"
	"fmt"
	"timeline/internal/model"
)

// Сохраняет пользователя и его почту/логин + пароль в хешированном виде
func (p *PostgresRepo) UserSave(ctx context.Context, user *model.UserInfo, creds *model.Credentials) (int, error) {
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
		INSERT INTO users (email, passwd_hash, user_name, telephone, social, about)
        VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING user_id;
	`
	userID := 0
	err = tx.QueryRowContext(ctx, query, creds.Login, string(creds.PasswdHash), user.Name, user.Telephone, user.Social, user.About).Scan(&userID)
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
func (p *PostgresRepo) UserByEmail(ctx context.Context, email string) (*model.User, *model.Credentials, error) {
	tx, err := p.db.Begin()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to start tx: %w", err)
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
	var user model.User
	var creds model.Credentials
	err = tx.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&creds.Login,
		&creds.PasswdHash,
		&user.Info.Name,
		&user.Info.Telephone,
		&user.Info.Social,
		&user.Info.About,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get user by email: %w", err)
	}
	err = tx.Commit()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to commit tx: %w", err)
	}
	return &user, &creds, nil
}

// Получить юзера по ID из БД
func (p *PostgresRepo) UserByID(ctx context.Context, id int) (*model.User, *model.Credentials, error) {
	tx, err := p.db.Begin()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to start tx: %w", err)
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
	var user model.User
	var creds model.Credentials
	err = tx.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&creds.Login,
		&creds.PasswdHash,
		&user.Info.Name,
		&user.Info.Telephone,
		&user.Info.Social,
		&user.Info.About,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to save user: %w", err)
	}
	err = tx.Commit()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to commit tx: %w", err)
	}
	return &user, &creds, nil
}

// Существуют ли юзер в БД
func (p *PostgresRepo) UserIsExist(ctx context.Context, email string) error {
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
        SELECT COUNT(*) FROM users
        WHERE email = $1;
    `

	var count int
	err = tx.QueryRowContext(ctx, query, email).Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to query user: %w", err)
	}

	// Если count == 0, значит, пользователь с таким email не существует
	if count == 0 {
		return ErrUserNotFound
	}
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit tx: %w", err)
	}
	return nil
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

	err = tx.QueryRowContext(ctx, query, code, user_id).Err()
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
func (p *PostgresRepo) UserCode(ctx context.Context, code string, user_id int) (string, error) {
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
		SELECT code FROM user_verify
        WHERE code = $1 AND user_id = $2;
	`
	var verifyCode string
	err = tx.QueryRowContext(ctx, query, code, user_id).Scan(&verifyCode)
	if err != nil {
		return "", fmt.Errorf("failed to save org code: %w", err)
	}
	err = tx.Commit()
	if err != nil {
		return "", fmt.Errorf("failed to commit tx: %w", err)
	}

	return verifyCode, nil
}
