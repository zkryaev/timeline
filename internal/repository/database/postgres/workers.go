package postgres

import (
	"context"
	"fmt"
	"timeline/internal/repository/models"
)

// Добавление работника к организации
func (p *PostgresRepo) WorkerAdd(ctx context.Context, worker *models.Worker) (int, error) {
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
		INSERT INTO workers
		(org_id, first_name, last_name, position, degree)
		VALUES($1, $2, $3, $4, $5)
		RETURNING worker_id;
	`
	var workerID int
	if err := tx.QueryRowContext(ctx, query,
		worker.OrgID,
		worker.FirstName,
		worker.LastName,
		worker.Position,
		worker.Degree,
	).Scan(&workerID); err != nil {
		return 0, fmt.Errorf("failed to add worker to org: %w", err)
	}
	if tx.Commit() != nil {
		return 0, fmt.Errorf("failed to commit transaction: %w", err)
	}
	return workerID, nil
}

func (p *PostgresRepo) Worker(ctx context.Context, WorkerID, OrgID int) (*models.Worker, error) {
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
		SELECT worker_id, org_id, first_name, last_name, position, degree
		FROM workers
		WHERE worker_id = $1
		AND org_id = $2;
	`
	var Worker models.Worker
	if err = tx.GetContext(ctx, &Worker, query, &WorkerID, &OrgID); err != nil {
		return nil, fmt.Errorf("failed to get worker: %w", err)
	}
	if tx.Commit() != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}
	return &Worker, nil
}

// Обновление информации работника организации
func (p *PostgresRepo) WorkerUpdate(ctx context.Context, worker *models.Worker) error {
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
		UPDATE workers
		SET
			first_name = $1,
			last_name = $2,
			position = $3,
			degree = $4
		WHERE worker_id = $5
		AND org_id = $6;
	`
	if err = tx.QueryRowContext(ctx, query,
		worker.FirstName,
		worker.LastName,
		worker.Position,
		worker.Degree,
		worker.WorkerID,
		worker.OrgID,
	).Err(); err != nil {
		return fmt.Errorf("failed to update worker: %w", err)
	}
	if tx.Commit() != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}

// Добавляет работника в предоставляемую услугу
func (p *PostgresRepo) WorkerAssignService(ctx context.Context, assignInfo *models.WorkerAssign) error {
	tx, err := p.db.Beginx()
	if err != nil {
		return fmt.Errorf("failed to start tx: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	// Проверка на существование и затем добавление в связующую таблицу
	query := `
		INSERT INTO worker_services
		(worker_id, service_id)
		SELECT 
			w.worker_id, s.service_id
		FROM workers w
		JOIN services s ON w.org_id = s.org_id
		WHERE w.worker_id = $1 
		AND s.service_id = $2
		AND s.org_id = $3;
	`
	rows, err := tx.ExecContext(ctx, query,
		assignInfo.WorkerID,
		assignInfo.ServiceID,
		assignInfo.OrgID,
	)
	if err != nil {
		return fmt.Errorf("failed to update worker: %w", err)
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

func (p *PostgresRepo) WorkerUnAssignService(ctx context.Context, assignInfo *models.WorkerAssign) error {
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
		DELETE FROM worker_services
		WHERE worker_id = (SELECT 
					worker_id
				FROM workers
				WHERE org_id = $1
				AND worker_id = $2)
		AND service_id = (SELECT 
					service_id
				FROM services
				WHERE org_id = $1
				AND service_id = $3);
	`
	rows, err := tx.ExecContext(ctx, query, &assignInfo.OrgID, &assignInfo.WorkerID, &assignInfo.ServiceID)
	if err != nil {
		return fmt.Errorf("failed to delete worker: %w", err)
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

// Получение списка работников организации
func (p *PostgresRepo) WorkerList(ctx context.Context, OrgID int) ([]*models.Worker, error) {
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
		SELECT worker_id, org_id, first_name, last_name, position, degree
		FROM workers
		WHERE org_id = $1;
	`
	Workers := make([]*models.Worker, 0, 3)
	if err = tx.SelectContext(ctx, &Workers, query, &OrgID); err != nil {
		return nil, fmt.Errorf("failed to get worker list: %w", err)
	}
	if tx.Commit() != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}
	return Workers, nil
}

// Удаляет работника из организации, а также связи с предоставляемыми услугами
func (p *PostgresRepo) WorkerDelete(ctx context.Context, WorkerID, OrgID int) error {
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
		DELETE FROM workers
		WHERE worker_id = $1
		AND org_id = $2;
	`
	rows, err := tx.ExecContext(ctx, query, &WorkerID, &OrgID)
	if err != nil {
		return fmt.Errorf("failed to delete worker: %w", err)
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
