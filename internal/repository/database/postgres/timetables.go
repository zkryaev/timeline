package postgres

import (
	"context"
	"fmt"
	"log"
	"timeline/internal/repository/models/orgmodel"
)

func (p *PostgresRepo) TimetableAdd(ctx context.Context, orgID int, new []*orgmodel.OpenHours) error {
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
			VALUES($1, $2, $3, $4, $5, $6)
	`
	for _, hours := range new {
		res, err := tx.ExecContext(ctx, query,
			hours.Weekday,
			hours.Open,
			hours.Close,
			hours.BreakStart,
			hours.BreakEnd,
			orgID,
		)
		if err != nil {
			return err
		}
		if res != nil {
			if _, errNoRowsAffected := res.RowsAffected(); errNoRowsAffected != nil {
				return ErrOrgNotFound
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
			AND ($2 <= 0 OR weekday = $2)
	`
	log.Println(orgID, weekday)
	res, err := tx.ExecContext(ctx, query, orgID, weekday)
	if err != nil {
		return err
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

// Принимает org_id и расписание организации. Если
func (p *PostgresRepo) TimetableUpdate(ctx context.Context, orgID int, new []*orgmodel.OpenHours) error {
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
	for _, hours := range new {
		res, err := tx.ExecContext(ctx, query,
			hours.Weekday,
			hours.Open,
			hours.Close,
			hours.BreakStart,
			hours.BreakEnd,
			orgID,
		)
		if err != nil {
			return err
		}
		if res != nil {
			if _, errNoRowsAffected := res.RowsAffected(); errNoRowsAffected != nil {
				return ErrOrgNotFound
			}
		}
	}
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit tx: %w", err)
	}
	return nil
}
