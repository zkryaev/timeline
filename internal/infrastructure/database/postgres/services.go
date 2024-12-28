package postgres

import (
	"context"
	"errors"
	"fmt"
	"timeline/internal/infrastructure/models/orgmodel"

	"github.com/jackc/pgx/v5"
)

var (
	ErrServiceNotFound = errors.New("service not found")
)

// Добавление организации предоставляемой услуги
func (p *PostgresRepo) ServiceAdd(ctx context.Context, service *orgmodel.Service) (int, error) {
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
		INSERT INTO services
		(org_id, name, cost, description)
		VALUES($1, $2, $3, $4)
		RETURNING service_id;
	`
	var serviceID int
	if err := tx.QueryRowContext(ctx, query,
		service.OrgID,
		service.Name,
		service.Cost,
		service.Description,
	).Scan(&serviceID); err != nil {
		return 0, err
	}
	if tx.Commit() != nil {
		return 0, fmt.Errorf("failed to commit transaction: %w", err)
	}
	return serviceID, nil
}

// Обновление информации о предоставляемой услуге
func (p *PostgresRepo) ServiceUpdate(ctx context.Context, service *orgmodel.Service) error {
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
		UPDATE services
		SET
			name = $1,
			cost = $2,
			description = $3
		WHERE is_delete = false
		AND service_id = $4 
		AND org_id = $5;
	`
	if err = tx.QueryRowContext(ctx, query,
		service.Name,
		service.Cost,
		service.Description,
		service.ServiceID,
		service.OrgID,
	).Err(); err != nil {
		return err
	}
	if tx.Commit() != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}

func (p *PostgresRepo) Service(ctx context.Context, ServiceID, OrgID int) (*orgmodel.Service, error) {
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
		SELECT service_id, org_id, name, cost, description
		FROM services
		WHERE is_delete = false 
		AND service_id = $1
		AND org_id = $2;
	`
	var Service orgmodel.Service
	if err = tx.GetContext(ctx, &Service, query, &ServiceID, &OrgID); err != nil {
		return nil, err
	}
	if tx.Commit() != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}
	return &Service, nil
}

// Получение списка услуг, предоставляемых организацией
func (p *PostgresRepo) ServiceList(ctx context.Context, OrgID int, Limit int, Offset int) ([]*orgmodel.Service, int, error) {
	tx, err := p.db.Beginx()
	if err != nil {
		return nil, 0, fmt.Errorf("failed to start tx: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	query := `SELECT 
			COUNT(*)
		FROM services 
		WHERE is_delete = false
		AND org_id = $1;
	`
	var found int
	if err = tx.QueryRowxContext(ctx, query, OrgID).Scan(&found); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, 0, ErrServiceNotFound
		}
		return nil, 0, fmt.Errorf("failed to get org's service list: %w", err)
	}
	query = `
		SELECT service_id, org_id, name, cost, description
		FROM services
		WHERE is_delete = false 
		AND org_id = $1
		LIMIT $2
		OFFSET $3;
	`
	Services := make([]*orgmodel.Service, 0, 3)
	if err = tx.SelectContext(ctx, &Services, query, &OrgID, Limit, Offset); err != nil {
		return nil, 0, err
	}
	if tx.Commit() != nil {
		return nil, 0, fmt.Errorf("failed to commit transaction: %w", err)
	}
	return Services, 0, nil
}

// Удаление услуги, предоставляемой организации и удаление связи с работниками
func (p *PostgresRepo) ServiceDelete(ctx context.Context, ServiceID, OrgID int) error {
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
		UPDATE services
		SET
			is_delete = true
		WHERE is_delete = false 
		AND service_id = $1
		AND org_id = $2;
	`
	rows, err := tx.ExecContext(ctx, query, &ServiceID, &OrgID)
	if err != nil {
		return err
	}
	if rows != nil {
		if _, NoAffectedRows := rows.RowsAffected(); NoAffectedRows != nil {
			return ErrServiceNotFound
		}
	}
	if tx.Commit() != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}

func (p *PostgresRepo) ServiceWorkerList(ctx context.Context, ServiceID, OrgID int) ([]*orgmodel.Worker, error) {
	tx, err := p.db.Beginx()
	if err != nil {
		return nil, fmt.Errorf("failed to start tx: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	query := `SELECT worker_id, org_id, first_name, last_name, position, degree
		FROM workers
		WHERE is_delete = false 
		AND worker_id IN 
		(SELECT
				worker_id 
			FROM worker_services 
			WHERE is_delete = false 
			AND service_id = $1
		)
		AND org_id = $2;
	`
	Workers := make([]*orgmodel.Worker, 0, 1)
	if err = tx.SelectContext(ctx, &Workers, query, &ServiceID, &OrgID); err != nil {
		return nil, err
	}
	if tx.Commit() != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}
	return Workers, nil
}
