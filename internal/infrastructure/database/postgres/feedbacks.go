package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"timeline/internal/infrastructure/models/recordmodel"

	"github.com/jackc/pgx/v5"
)

var (
	ErrFeedbackNotFound = errors.New("feedback not found")
)

func (p *PostgresRepo) FeedbackList(ctx context.Context, params *recordmodel.FeedbackParams) ([]*recordmodel.Feedback, int, error) {
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
		FROM records r
		JOIN feedbacks f
		ON r.record_id = f.record_id AND reviewed = true
		WHERE ($1 <= 0 OR f.record_id = $1)
		AND ($2 <= 0 OR r.user_id = $2)
		AND ($3 <= 0 OR r.org_id = $3);
	`
	var found int
	if err = tx.QueryRowxContext(ctx, query, params.RecordID, params.UserID, params.OrgID).Scan(&found); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, 0, ErrServiceNotFound
		}
		return nil, 0, fmt.Errorf("failed to get org's service list: %w", err)
	}
	query = `SELECT f.record_id, f.stars, f.feedback
		FROM feedbacks f
		JOIN records r
		ON r.record_id = f.record_id AND r.reviewed = true
		WHERE ($1 <= 0 OR f.record_id = $1)
		AND ($2 <= 0 OR r.user_id = $2)
		AND ($3 <= 0 OR r.org_id = $3);
	`
	feedbacks := make([]*recordmodel.Feedback, 0, 1)
	if err := tx.SelectContext(ctx, &feedbacks, query, params.RecordID, params.UserID, params.OrgID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, 0, ErrFeedbackNotFound
		}
		return nil, 0, err
	}
	if err = tx.Commit(); err != nil {
		return nil, 0, fmt.Errorf("failed to commit tx: %w", err)
	}
	return feedbacks, 0, nil
}

func (p *PostgresRepo) FeedbackSet(ctx context.Context, feedback *recordmodel.Feedback) error {
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
		INSERT INTO feedbacks
		(record_id, stars, feedback)
		SELECT $1, $2, $3
		FROM slots s
		WHERE s.record_id = $1
		AND CURRENT_TIMESTAMP >= s.session_end;
	`
	rows, err := tx.ExecContext(
		ctx,
		query,
		feedback.RecordID,
		feedback.Stars,
		feedback.Feedback,
	)
	if err != nil {
		return err
	}
	if rows != nil {
		if rowsAffected, _ := rows.RowsAffected(); rowsAffected == 0 {
			return fmt.Errorf("no rows inserted")
		}
	}
	query = `UPDATE records
		SET
			reviewed = true
		WHERE record_id = $1;
	`
	rows, err = tx.ExecContext(
		ctx,
		query,
		feedback.RecordID,
	)
	if err != nil {
		return err
	}
	if rows != nil {
		if rowsAffected, _ := rows.RowsAffected(); rowsAffected == 0 {
			return fmt.Errorf("no rows inserted")
		}
	}
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit tx: %w", err)
	}
	return nil
}

func (p *PostgresRepo) FeedbackUpdate(ctx context.Context, feedback *recordmodel.Feedback) error {
	tx, err := p.db.Beginx()
	if err != nil {
		return fmt.Errorf("failed to start tx: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	query := `UPDATE feedbacks
		SET
			stars = $1,
			feedback = $2
		WHERE record_id = $3
	`
	rows, err := tx.ExecContext(
		ctx,
		query,
		feedback.Stars,
		feedback.Feedback,
		feedback.RecordID,
	)
	if err != nil {
		return err
	}
	if rows != nil {
		if rowsAffected, _ := rows.RowsAffected(); rowsAffected == 0 {
			return fmt.Errorf("no rows inserted")
		}
	}
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit tx: %w", err)
	}
	return nil
}

func (p *PostgresRepo) FeedbackDelete(ctx context.Context, params *recordmodel.FeedbackParams) error {
	tx, err := p.db.Beginx()
	if err != nil {
		return fmt.Errorf("failed to start tx: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	query := `DELETE FROM feedbacks
		WHERE record_id = $1
	`
	rows, err := tx.ExecContext(ctx, query, params.RecordID)
	if err != nil {
		return err
	}
	if rows != nil {
		if rowsAffected, _ := rows.RowsAffected(); rowsAffected == 0 {
			return fmt.Errorf("no rows inserted")
		}
	}
	query = `UPDATE records
		SET
			reviewed = false
		WHERE record_id = $1
	`
	rows, err = tx.ExecContext(
		ctx,
		query,
		params.RecordID,
	)
	if err != nil {
		return err
	}
	if rows != nil {
		if rowsAffected, _ := rows.RowsAffected(); rowsAffected == 0 {
			return fmt.Errorf("no rows inserted")
		}
	}
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit tx: %w", err)
	}
	return nil
}
