package postgres

import (
	"context"
	"errors"
	"fmt"
	"timeline/internal/infrastructure/models/orgmodel"

	"github.com/jackc/pgx/v5"
)

var (
	ErrWorkerNotFound = errors.New("worker not found")
)

// Добавление работника к организации
func (p *PostgresRepo) WorkerAdd(ctx context.Context, worker *orgmodel.Worker) (int, error) {
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
		(org_id, first_name, last_name, position, degree, session_duration)
		VALUES($1, $2, $3, $4, $5, $6)
		RETURNING worker_id;
	`
	var workerID int
	if err := tx.QueryRowContext(ctx, query,
		worker.OrgID,
		worker.FirstName,
		worker.LastName,
		worker.Position,
		worker.Degree,
		worker.SessionDuration,
	).Scan(&workerID); err != nil {
		return 0, fmt.Errorf("failed to add worker to org: %w", err)
	}
	if tx.Commit() != nil {
		return 0, fmt.Errorf("failed to commit transaction: %w", err)
	}
	return workerID, nil
}

func (p *PostgresRepo) Worker(ctx context.Context, WorkerID, OrgID int) (*orgmodel.Worker, error) {
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
		SELECT worker_id, COALESCE(uuid, '') AS uuid, org_id, first_name, last_name, position, degree, session_duration
		FROM workers
		WHERE is_delete = false
		AND worker_id = $1
		AND org_id = $2;
	`
	var Worker orgmodel.Worker
	if err = tx.GetContext(ctx, &Worker, query, &WorkerID, &OrgID); err != nil {
		return nil, fmt.Errorf("failed to get worker: %w", err)
	}
	if tx.Commit() != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}
	return &Worker, nil
}

// Обновление информации работника организации
func (p *PostgresRepo) WorkerUpdate(ctx context.Context, worker *orgmodel.Worker) error {
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
			degree = $4,
			session_duration = $5
		WHERE is_delete = false 
		AND worker_id = $6
		AND org_id = $7;
	`
	if err = tx.QueryRowContext(ctx, query,
		worker.FirstName,
		worker.LastName,
		worker.Position,
		worker.Degree,
		worker.SessionDuration,
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

// Частичное обновление информации работника организации
func (p *PostgresRepo) WorkerPatch(ctx context.Context, worker *orgmodel.Worker) error {
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
			first_name = COALESCE(NULLIF($1, ''), first_name),
			last_name = COALESCE(NULLIF($2, ''), last_name),
			position = COALESCE(NULLIF($3, ''), position),
			degree = COALESCE(NULLIF($4, ''), degree),
			session_duration = COALESCE(NULLIF($5, 0), session_duration)
		WHERE is_delete = false
		AND worker_id = $6
		AND org_id = $7;
	`
	if err = tx.QueryRowContext(ctx, query,
		worker.FirstName,
		worker.LastName,
		worker.Position,
		worker.Degree,
		worker.SessionDuration,
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
func (p *PostgresRepo) WorkerAssignService(ctx context.Context, assignInfo *orgmodel.WorkerAssign) error {
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
		WHERE w.is_delete = false
		AND s.is_delete = false
		AND w.worker_id = $1 
		AND s.service_id = $2
		AND s.org_id = $3;
	`
	res, err := tx.ExecContext(ctx, query,
		assignInfo.WorkerID,
		assignInfo.ServiceID,
		assignInfo.OrgID,
	)
	switch {
	case err != nil:
		return fmt.Errorf("failed to update worker: %w", err)
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

func (p *PostgresRepo) WorkerUnAssignService(ctx context.Context, assignInfo *orgmodel.WorkerAssign) error {
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
		UPDATE worker_services
		SET
			is_delete = true
		WHERE worker_id = (SELECT 
									worker_id
								FROM workers
								WHERE is_delete = false 
								AND org_id = $1
								AND worker_id = $2
							)
		AND service_id = (SELECT 
									service_id
								FROM services
								WHERE is_delete = false 
								AND org_id = $1
								AND service_id = $3
							);
	`
	res, err := tx.ExecContext(ctx, query, &assignInfo.OrgID, &assignInfo.WorkerID, &assignInfo.ServiceID)
	switch {
	case err != nil:
		return fmt.Errorf("failed to unassign worker from service: %w", err)
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

// Получение списка работников организации
func (p *PostgresRepo) WorkerList(ctx context.Context, OrgID, Limit, Offset int) ([]*orgmodel.Worker, int, error) {
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
			COUNT(worker_id)
		FROM workers 
		WHERE is_delete = false
		AND org_id = $1;
	`
	var found int
	if err = tx.QueryRowxContext(ctx, query, OrgID).Scan(&found); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, 0, ErrWorkerNotFound
		}
		return nil, 0, fmt.Errorf("failed to get org's service list: %w", err)
	}
	query = `
		SELECT worker_id, COALESCE(uuid, '') AS uuid, org_id, first_name, last_name, position, degree, session_duration
		FROM workers
		WHERE is_delete = false 
		AND org_id = $1
		LIMIT $2
		OFFSET $3;
	`
	Workers := make([]*orgmodel.Worker, 0, 3)
	if err = tx.SelectContext(ctx, &Workers, query, OrgID, Limit, Offset); err != nil {
		return nil, 0, fmt.Errorf("failed to get worker list: %w", err)
	}
	if tx.Commit() != nil {
		return nil, 0, fmt.Errorf("failed to commit transaction: %w", err)
	}
	return Workers, found, nil
}

func (p *PostgresRepo) WorkerSoftDelete(ctx context.Context, WorkerID, OrgID int) error {
	tx, err := p.db.Beginx()
	if err != nil {
		return fmt.Errorf("failed to start tx: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	mainQuery := `
		UPDATE workers
		SET
			is_delete = TRUE
		WHERE is_delete = FALSE
		AND worker_id = $1
		AND org_id = $2;
	`
	res, err := tx.ExecContext(ctx, mainQuery, WorkerID, OrgID)
	switch {
	case err != nil:
		return fmt.Errorf("failed to delete worker: %w", err)
	default:
		if rowsAffected, _ := res.RowsAffected(); rowsAffected == 0 {
			return ErrNoRowsAffected
		}
	}
	queries := []string{
		`
		UPDATE worker_schedules
		SET
			is_delete = TRUE
		WHERE worker_id = $1
		AND is_delete = FALSE;
	`,
		`
		UPDATE worker_services
		SET
			is_delete = TRUE
		WHERE worker_id = $1
		AND is_delete = FALSE;
	`,
	}
	for _, query := range queries {
		if _, err := tx.ExecContext(ctx, query, WorkerID); err != nil {
			return fmt.Errorf("failed to delete worker's deps: %w", err)
		}
	}
	if tx.Commit() != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}

// Only for tests!
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
	res, err := tx.ExecContext(ctx, query, &WorkerID, &OrgID)
	switch {
	case err != nil:
		return fmt.Errorf("failed to delete worker: %w", err)
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

func (p *PostgresRepo) WorkerUUID(ctx context.Context, workerID int) (string, error) {
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
		SELECT uuid
		FROM workers
		WHERE worker_id = $1;
	`
	var uuid string
	if err := tx.QueryRowContext(ctx, query, workerID).Scan(&uuid); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", ErrOrgNotFound
		}
		return "", fmt.Errorf("failed to get worker uuid by id: %w", err)
	}
	if err = tx.Commit(); err != nil {
		return "", fmt.Errorf("failed to commit tx: %w", err)
	}
	return uuid, nil
}
func (p *PostgresRepo) WorkerSetUUID(ctx context.Context, workerID int, NewUUID string) error {
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
			uuid = $1
		WHERE worker_id = $2;
	`
	res, err := tx.ExecContext(ctx, query, NewUUID, workerID)
	switch {
	case err != nil:
		return fmt.Errorf("failed to set worker's uuid: %w", err)
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

func (p *PostgresRepo) WorkerDeleteURL(ctx context.Context, URL string) error {
	tx, err := p.db.Beginx()
	if err != nil {
		return fmt.Errorf("failed to start tx: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	query := `UPDATE workers SET uuid = '' WHERE uuid = $1;`
	res, err := tx.ExecContext(ctx, query, URL)
	switch {
	case err != nil:
		return fmt.Errorf("failed to delete worker's url: %w", err)
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
