package postgres

import (
	"context"
	"errors"
	"fmt"
	"timeline/internal/repository/models/orgmodel"

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
			COUNT(*)
		FROM worker_schedules
		WHERE org_id = $1;
	`
	var found int
	if err = tx.QueryRowxContext(ctx, query, metainfo.OrgID).Scan(&found); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrScheduleNotFound
		}
		return nil, fmt.Errorf("failed to get org's service list: %w", err)
	}
	query = `
		SELECT 
			worker_schedule_id,
			weekday, 
			start, 
			over
		FROM worker_schedules
        WHERE ($1 <= 0 OR worker_id = $1) 
		AND org_id = $2
		AND ($3 <= 0 OR weekday = $3)
		LIMIT $4
		OFFSET $5;
	`
	schedule := make([]*orgmodel.Schedule, 0, 7) // 7 = число дней в неделе
	if err = tx.SelectContext(
		ctx,
		&schedule,
		query,
		metainfo.WorkerID,
		metainfo.OrgID,
		metainfo.Weekday,
		metainfo.Limit,
		metainfo.Offset,
	); err != nil {
		return nil, err
	}
	scheduleList := &orgmodel.ScheduleList{
		WorkerID: metainfo.WorkerID,
		OrgID:    metainfo.OrgID,
		Schedule: schedule,
		Found:    found,
	}
	if tx.Commit() != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}
	return scheduleList, nil
}

// Добавить для указанного рабонтика расписание на 1 день - неделю
// Обновляет длительность сеанса работника
func (p *PostgresRepo) AddWorkerSchedule(ctx context.Context, schedule *orgmodel.ScheduleList) error {
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
			SELECT open::time AS open_time, close::time AS close_time
			FROM timetables
			WHERE weekday = $1 AND org_id = $4
		)
		INSERT INTO worker_schedules (weekday, start, over, org_id, worker_id)
		SELECT $1, $2, $3, $4, $5
		FROM orgschedule
		WHERE
			$6::time >= orgschedule.open_time 
		AND $6::time <= orgschedule.close_time
		AND $7::time >= orgschedule.open_time 
		AND $7::time <= orgschedule.close_time;
	`
	for _, v := range schedule.Schedule {
		rows, err := tx.ExecContext(ctx, query, v.Weekday, v.Start, v.Over, schedule.OrgID, schedule.WorkerID, v.Start, v.Over)
		if err != nil {
			return err
		}
		if rows != nil {
			if rowsAffected, _ := rows.RowsAffected(); rowsAffected == 0 {
				return fmt.Errorf("no rows inserted")
			}
		}
	}
	query = `
		UPDATE workers
		SET
			session_duration = COALESCE(NULLIF($1, 0), session_duration)
		WHERE
			worker_id = $2
	`
	rows, err := tx.ExecContext(ctx, query, schedule.SessionDuration, schedule.WorkerID)
	if err != nil {
		return err
	}
	if rows != nil {
		if rowsAffected, _ := rows.RowsAffected(); rowsAffected == 0 {
			return fmt.Errorf("no rows inserted")
		}
	}
	if tx.Commit() != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}

// Обновление расписания работника на всю неделю
func (p *PostgresRepo) UpdateWorkerSchedule(ctx context.Context, schedule *orgmodel.ScheduleList) error {
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
		WHERE worker_schedule_id = $6
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
		rows, err := tx.ExecContext(ctx, query, v.Weekday, v.Start, v.Over, schedule.OrgID, schedule.WorkerID, v.WorkerScheduleID, v.Start, v.Over)
		if err != nil {
			return err
		}
		if rows != nil {
			if rowsAffected, _ := rows.RowsAffected(); rowsAffected == 0 {
				return fmt.Errorf("no rows inserted")
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
		AND ($3 <= 0 OR weekday = $3)
	`
	rows, err := tx.ExecContext(ctx, query, metainfo.WorkerID, metainfo.OrgID, metainfo.Weekday)
	if err != nil {
		return err
	}
	if rows != nil {
		if rowsAffected, _ := rows.RowsAffected(); rowsAffected == 0 {
			return ErrScheduleNotFound
		}
	}
	if tx.Commit() != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}
