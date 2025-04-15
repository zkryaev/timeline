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
			COUNT(r.record_id)
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
	query = `
		SELECT 
			f.record_id, 
			f.stars, 
			f.feedback, 
			s.name AS service_name, 
			w.first_name AS worker_first_name, 
			w.last_name worker_last_name,
			u.first_name AS user_first_name,
			u.last_name AS user_last_name,
			r.created_at AS record_date
		FROM feedbacks f
		JOIN records r ON r.record_id = f.record_id AND r.reviewed = true
		JOIN users u ON r.user_id = u.user_id
        JOIN workers w ON r.worker_id = w.worker_id
        JOIN services s ON r.service_id = s.service_id
		WHERE ($1 <= 0 OR f.record_id = $1)
		AND ($2 <= 0 OR r.user_id = $2)
		AND ($3 <= 0 OR r.org_id = $3)
		ORDER BY f.created_at DESC
		LIMIT $4
		OFFSET $5;
	`
	feedbacks := make([]*recordmodel.Feedback, 0, 1)
	if err = tx.SelectContext(ctx, &feedbacks, query, params.RecordID, params.UserID, params.OrgID, params.Limit, params.Offset); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, 0, ErrFeedbackNotFound
		}
		return nil, 0, err
	}
	if err = tx.Commit(); err != nil {
		return nil, 0, fmt.Errorf("failed to commit tx: %w", err)
	}
	return feedbacks, found, nil
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
		FROM records r
		JOIN slots s ON s.slot_id = r.slot_id 
		WHERE r.record_id = $1 AND r.user_id = $4
		AND CURRENT_TIMESTAMP >= s.session_end;
	`
	res, err := tx.ExecContext(
		ctx,
		query,
		feedback.RecordID,
		feedback.Stars,
		feedback.Feedback,
		feedback.TData.ID,
	)
	switch {
	case err != nil:
		return fmt.Errorf("failed to set feedback: %w", err)
	default:
		if rowsAffected, _ := res.RowsAffected(); rowsAffected == 0 {
			return ErrNoRowsAffected
		}
	}
	query = `UPDATE records
		SET
			reviewed = true
		WHERE record_id = $1;
	`
	res, err = tx.ExecContext(
		ctx,
		query,
		feedback.RecordID,
	)
	switch {
	case err != nil:
		return fmt.Errorf("failed to set feedback: %w", err)
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
	query := `UPDATE feedbacks f
		SET
			stars = $1,
			feedback = $2
		FROM records r
		WHERE f.record_id = r.record_id 
		AND r.record_id = $3
		AND r.user_id = $4;
	`
	res, err := tx.ExecContext(
		ctx,
		query,
		feedback.Stars,
		feedback.Feedback,
		feedback.RecordID,
		feedback.TData.ID,
	)
	switch {
	case err != nil:
		return fmt.Errorf("failed to update feedback: %w", err)
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
		WHERE record_id = $1;
	`
	res, err := tx.ExecContext(ctx, query, params.RecordID)
	switch {
	case err != nil:
		return fmt.Errorf("failed to delete feedback: %w", err)
	default:
		if rowsAffected, _ := res.RowsAffected(); rowsAffected == 0 {
			return ErrNoRowsAffected
		}
	}
	query = `UPDATE records
		SET
			reviewed = false
		WHERE record_id = $1;
	`
	res, err = tx.ExecContext(
		ctx,
		query,
		params.RecordID,
	)
	switch {
	case err != nil:
		return fmt.Errorf("failed to delete feedback: %w", err)
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
