package postgres

import (
	"context"
	"errors"
	"fmt"
	"timeline/internal/infrastructure/models"
	"timeline/internal/infrastructure/models/orgmodel"

	"github.com/jackc/pgx/v5"
)

var (
	ErrTimetableNotFound = errors.New("timetable not found")
)

func (p *PostgresRepo) Timetable(ctx context.Context, orgID int, tdata models.TokenData) ([]*orgmodel.OpenHours, string, error) {
	tx, err := p.db.Beginx()
	if err != nil {
		return nil, "", fmt.Errorf("failed to start tx: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	query := `
		SELECT weekday, open, close, break_start, break_end
		FROM timetables
		WHERE org_id = $1;
	`
	timetable := make([]*orgmodel.OpenHours, 0, 1)
	if err = tx.SelectContext(ctx, &timetable, query, orgID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, "", ErrTimetableNotFound
		}
		return nil, "", err
	}
	var city string
	switch tdata.IsOrg {
	case false:
		query = `
			SELECT city
			FROM users
			WHERE user_id = $1;
		`
	case true:
		query = `
			SELECT city
			FROM orgs
			WHERE org_id = $1;
		`
	}
	if err := tx.QueryRowxContext(ctx, query, tdata.ID).Scan(&city); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, "", ErrTimetableNotFound
		}
		return nil, "", err
	}
	if err = tx.Commit(); err != nil {
		return nil, "", fmt.Errorf("failed to commit tx: %w", err)
	}
	return timetable, city, nil
}

func (p *PostgresRepo) TimetableAdd(ctx context.Context, orgID int, newTime []*orgmodel.OpenHours) error {
	tx, err := p.db.Beginx()
	if err != nil {
		return fmt.Errorf("failed to start tx: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	query := `INSERT INTO timetables
			(weekday, open, close, break_start, break_end, org_id)
			VALUES($1, $2, $3, $4, $5, $6);
	`
	for _, hours := range newTime {
		res, err := tx.ExecContext(ctx, query, //nolint:govet // ...
			hours.Weekday,
			hours.Open,
			hours.Close,
			hours.BreakStart,
			hours.BreakEnd,
			orgID,
		)
		switch {
		case err != nil:
			return fmt.Errorf("failed to add timetable: %w", err)
		default:
			if rowsAffected, _ := res.RowsAffected(); rowsAffected == 0 {
				return ErrNoRowsAffected
			}
		}
	}
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit tx: %w", err)
	}
	return nil
}

func (p *PostgresRepo) TimetableDelete(ctx context.Context, orgID, weekday int) error {
	tx, err := p.db.Beginx()
	if err != nil {
		return fmt.Errorf("failed to start tx: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	query := `DELETE FROM timetables
			WHERE org_id = $1
			AND ($2 <= 0 OR weekday = $2);
	`
	res, err := tx.ExecContext(ctx, query, orgID, weekday)
	switch {
	case err != nil:
		return fmt.Errorf("failed to delete timetable: %w", err)
	default:
		if rowsAffected, _ := res.RowsAffected(); rowsAffected == 0 {
			return ErrNoRowsAffected
		}
	}
	query = `
		UPDATE worker_schedules
		SET 
			is_delete = true
		WHERE is_delete = false
		AND org_id = $1
		AND ($2 <= 0 OR weekday = $2);
	`
	if _, err = tx.ExecContext(ctx, query, orgID, weekday); err != nil {
		return fmt.Errorf("failed to delete timetable: %w", err)
	}
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit tx: %w", err)
	}
	return nil
}

// Принимает org_id и расписание организации. Если
func (p *PostgresRepo) TimetableUpdate(ctx context.Context, orgID int, timeList []*orgmodel.OpenHours) error {
	tx, err := p.db.Beginx()
	if err != nil {
		return fmt.Errorf("failed to start tx: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	query := `UPDATE timetables
		SET 
			weekday = $1,
			open = $2,
			close = $3,
			break_start = $4,
			break_end = $5
		WHERE org_id = $6
		AND weekday = $1;
	`
	for _, hours := range timeList {
		res, err := tx.ExecContext(ctx, query, //nolint:govet // ...
			hours.Weekday,
			hours.Open,
			hours.Close,
			hours.BreakStart,
			hours.BreakEnd,
			orgID,
		)
		switch {
		case err != nil:
			return fmt.Errorf("failed to update timetable: %w", err)
		default:
			if rowsAffected, _ := res.RowsAffected(); rowsAffected == 0 {
				return ErrNoRowsAffected
			}
		}
	}
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit tx: %w", err)
	}
	return nil
}
