package postgres

import (
	"context"
	"errors"
	"fmt"
	"timeline/internal/infrastructure/models/orgmodel"

	"github.com/jackc/pgx/v5"
)

var (
	ErrScheduleNotFound = errors.New("schedule not found")
)

// Получение расписания работника.
// Если weekday = 0, то получаем расписание на всю неделю
// Иначе заданный день. (1.Пн...7.Вс)
func (p *PostgresRepo) WorkerSchedule(ctx context.Context, metainfo *orgmodel.ScheduleParams) (*orgmodel.ScheduleList, error) {
	tx, err := p.db.Beginx()
	if err != nil {
		return nil, fmt.Errorf("failed to start tx: %w", err)
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
		AND ($1 <= 0 OR org_id = $1)
		AND org_id = $1
		AND ($2 <= 0 OR worker_id = $2);
	`
	var found int
	if err = tx.QueryRowxContext(ctx, query, metainfo.OrgID, metainfo.WorkerID).Scan(&found); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrScheduleNotFound
		}
		return nil, fmt.Errorf("failed to get org's service list: %w", err)
	}
	query = `
		SELECT 
			worker_id, session_duration
		FROM workers
        WHERE is_delete = false 
		AND ($1 <= 0 OR worker_id = $1) 
		AND org_id = $2
		LIMIT $3
		OFFSET $4;
	`
	workerList := make([]struct {
		WorkerID        int `db:"worker_id"`
		SessionDuration int `db:"session_duration"`
	}, 0, 3)
	if err = tx.SelectContext(
		ctx,
		&workerList,
		query,
		metainfo.WorkerID,
		metainfo.OrgID,
		metainfo.Limit,
		metainfo.Offset,
	); err != nil {
		return nil, err
	}
	// Получение для каждого воркера его расписания
	query = `
		SELECT 
			worker_schedule_id,
			weekday, 
			start, 
			over
		FROM worker_schedules
        WHERE is_delete = false
		AND worker_id = $1 
		AND org_id = $2
		AND ($3 <= 0 OR weekday = $3);
	`
	resp := &orgmodel.ScheduleList{
		Workers: make([]*orgmodel.WorkerSchedule, 0, len(workerList)),
		Found:   found,
	}
	var schedule []*orgmodel.Schedule
	for _, worker := range workerList {
		schedule = make([]*orgmodel.Schedule, 0, 7) // 7 - дней в неделе
		if err = tx.SelectContext(
			ctx,
			&schedule,
			query,
			worker.WorkerID,
			metainfo.OrgID,
			metainfo.Weekday,
		); err != nil {
			return nil, err
		}
		worker := &orgmodel.WorkerSchedule{
			WorkerID:        worker.WorkerID,
			OrgID:           metainfo.OrgID,
			SessionDuration: worker.SessionDuration,
			Schedule:        schedule,
		}
		resp.Workers = append(resp.Workers, worker)
	}
	if tx.Commit() != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}
	return resp, nil
}

// Добавить для указанного рабонтика расписание на 1 день - неделю
// Обновляет длительность сеанса работника
func (p *PostgresRepo) AddWorkerSchedule(ctx context.Context, schedule *orgmodel.WorkerSchedule) error {
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
		WITH orgschedule AS (
			SELECT 
				org_id,
				open::time AS open_time, 
				close::time AS close_time, 
				break_start::time AS break_start, 
				break_end::time AS break_end
			FROM timetables
			WHERE weekday = $1 
			AND org_id = $4
		)
		INSERT INTO worker_schedules 
			(weekday, start, over, org_id, worker_id)
		SELECT $1, $2, $3, $4, $5
		FROM (SELECT 1) AS dummy_table
		WHERE 
		NOT EXISTS (
			SELECT 1
			FROM worker_schedules
			WHERE is_delete = false
			AND worker_id = $5
			AND weekday = $1
		)
		AND EXISTS (
			SELECT 1
			FROM orgschedule
			WHERE $6::time >= open_time
			AND $6::time <= close_time
			AND $7::time >= open_time 
			AND $7::time <= close_time
			AND (
				$7::time <= break_start OR $6::time >= break_end OR 
				($6::time < break_start AND $7::time > break_end)
			) 
		);
	`
	for _, v := range schedule.Schedule {
		res, err := tx.ExecContext(ctx, query, v.Weekday, v.Start, v.Over, schedule.OrgID, schedule.WorkerID, v.Start, v.Over) //nolint:govet // ...
		switch {
		case err != nil:
			return fmt.Errorf("failed to add worker schedule: %w", err)
		default:
			if rowsAffected, _ := res.RowsAffected(); rowsAffected == 0 {
				return ErrNoRowsAffected
			}
		}
	}
	query = `
		UPDATE workers
		SET
			session_duration = COALESCE(NULLIF($1, 0), session_duration)
		WHERE is_delete = false
		AND worker_id = $2;
	`
	res, err := tx.ExecContext(ctx, query, schedule.SessionDuration, schedule.WorkerID)
	switch {
	case err != nil:
		return fmt.Errorf("failed to add worker schedule: %w", err)
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

// Обновление расписания работника на всю неделю
func (p *PostgresRepo) UpdateWorkerSchedule(ctx context.Context, schedule *orgmodel.WorkerSchedule) error {
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
			session_duration = $1
		WHERE is_delete = false
		AND worker_id = $2
		AND org_id = $3;
	`
	res, err := tx.ExecContext(ctx, query, schedule.SessionDuration, schedule.WorkerID, schedule.OrgID) //nolint:govet // ...
	switch {
	case err != nil:
		return fmt.Errorf("failed to update worker schedule: %w", err)
	default:
		if rowsAffected, _ := res.RowsAffected(); rowsAffected == 0 {
			return ErrNoRowsAffected
		}
	}
	query = `
		WITH orgschedule AS (
			SELECT open::time AS open_time, close::time AS close_time
			FROM timetables
			WHERE weekday = $1 AND org_id = $4
		)
		UPDATE worker_schedules
		SET weekday = $1, 
			start = $2, 
			over = $3, 
			org_id = $4, 
			worker_id = $5
		WHERE is_delete = false
		AND worker_schedule_id = $6
		AND EXISTS (
			SELECT 1
			FROM orgschedule 
			WHERE $7::time >= orgschedule.open_time
			AND $7::time <= orgschedule.close_time
		)
		AND EXISTS (
			SELECT 1
			FROM orgschedule 
			WHERE $8::time >= orgschedule.open_time
			AND $8::time <= orgschedule.close_time
		);
	`
	for _, v := range schedule.Schedule {
		res, err := tx.ExecContext(ctx, query, v.Weekday, v.Start, v.Over, schedule.OrgID, schedule.WorkerID, v.WorkerScheduleID, v.Start, v.Over) //nolint:govet // ...
		switch {
		case err != nil:
			return fmt.Errorf("failed to update worker schedule: %w", err)
		default:
			if rowsAffected, _ := res.RowsAffected(); rowsAffected == 0 {
				return ErrNoRowsAffected
			}
		}
	}
	if tx.Commit() != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}

// Если weekday = 0, то удаляется расписание на всю неделю
// Иначе заданный день. (1.Пн...7.Вс)
func (p *PostgresRepo) SoftDeleteWorkerSchedule(ctx context.Context, metainfo *orgmodel.ScheduleParams) error {
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
		UPDATE worker_schedules
		SET 
			is_delete = true
		WHERE is_delete = false
		AND worker_id = $1
		AND org_id = $2
		AND ($3 <= 0 OR weekday = $3);
	`
	res, err := tx.ExecContext(ctx, query, metainfo.WorkerID, metainfo.OrgID, metainfo.Weekday)
	switch {
	case err != nil:
		return fmt.Errorf("failed to delete worker schedule: %w", err)
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

// Only for tests!
func (p *PostgresRepo) DeleteWorkerSchedule(ctx context.Context, metainfo *orgmodel.ScheduleParams) error {
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
		DELETE FROM worker_schedules
		WHERE worker_id = $1
		AND org_id = $2
		AND ($3 <= 0 OR weekday = $3);
	`
	res, err := tx.ExecContext(ctx, query, metainfo.WorkerID, metainfo.OrgID, metainfo.Weekday)
	switch {
	case err != nil:
		return fmt.Errorf("failed to delete worker schedule: %w", err)
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
