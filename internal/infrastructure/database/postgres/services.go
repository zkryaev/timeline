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
	if err = tx.QueryRowContext(ctx, query,
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

func (p *PostgresRepo) Service(ctx context.Context, serviceID, orgID int) (*orgmodel.Service, error) {
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
	var service orgmodel.Service
	if err = tx.GetContext(ctx, &service, query, &serviceID, &orgID); err != nil {
		return nil, err
	}
	if tx.Commit() != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}
	return &service, nil
}

// Получение списка услуг, предоставляемых организацией
func (p *PostgresRepo) ServiceList(ctx context.Context, orgID int, limit int, offset int) ([]*orgmodel.Service, int, error) {
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
			COUNT(service_id)
		FROM services 
		WHERE is_delete = false
		AND org_id = $1;
	`
	var found int
	if err = tx.QueryRowxContext(ctx, query, orgID).Scan(&found); err != nil {
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
	services := make([]*orgmodel.Service, 0, 3)
	if err = tx.SelectContext(ctx, &services, query, &orgID, limit, offset); err != nil {
		return nil, 0, err
	}
	if tx.Commit() != nil {
		return nil, 0, fmt.Errorf("failed to commit transaction: %w", err)
	}
	return services, found, nil
}

func (p *PostgresRepo) ServiceSoftDelete(ctx context.Context, serviceID, orgID int) error {
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
			is_delete = TRUE
		WHERE is_delete = FALSE 
		AND service_id = $1
		AND org_id = $2;
	`
	res, err := tx.ExecContext(ctx, query, &serviceID, &orgID)
	switch {
	case err != nil:
		return fmt.Errorf("failed to delete selected service: %w", err)
	default:
		if rowsAffected, _ := res.RowsAffected(); rowsAffected == 0 {
			return ErrNoRowsAffected
		}
	}
	query = `
		UPDATE worker_services
		SET
			is_delete = TRUE
		WHERE service_id = $1;
	`
	if _, err = tx.ExecContext(ctx, query, &serviceID); err != nil {
		return fmt.Errorf("failed to delete selected service: %w", err)
	}
	if tx.Commit() != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}

func (p *PostgresRepo) ServiceDelete(ctx context.Context, serviceID, orgID int) error {
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
		DELETE FROM services
		WHERE service_id = $1
		AND org_id = $2;
	`
	res, err := tx.ExecContext(ctx, query, &serviceID, &orgID)
	switch {
	case err != nil:
		return fmt.Errorf("failed to delete selected service: %w", err)
	default:
		if rowsAffected, _ := res.RowsAffected(); rowsAffected == 0 {
			return ErrNoRowsAffected
		}
	}
	if tx.Commit() != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}

func (p *PostgresRepo) ServiceWorkerList(ctx context.Context, serviceID, orgID int) ([]*orgmodel.Worker, error) {
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
	workers := make([]*orgmodel.Worker, 0, 1)
	if err = tx.SelectContext(ctx, &workers, query, &serviceID, &orgID); err != nil {
		return nil, err
	}
	if tx.Commit() != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}
	return workers, nil
}
