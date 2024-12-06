package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"timeline/internal/repository/models/recordmodel"
)

var (
	ErrFeedbackNotFound = errors.New("feedback not found")
)

func (p *PostgresRepo) Feedback(ctx context.Context, params *recordmodel.FeedbackParams) (*recordmodel.Feedback, error) {
	tx, err := p.db.Beginx()
	if err != nil {
		return nil, fmt.Errorf("failed to start tx: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	query := `SELECT record_id, stars, feedback
		FROM feedbacks
		WHERE record_id = $1		
	`
	feedback := &recordmodel.Feedback{}
	if err := tx.GetContext(ctx, feedback, query, params.RecordID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrFeedbackNotFound
		}
		return nil, err
	}
	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit tx: %w", err)
	}
	return feedback, nil
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
	query := `INSERT INTO feedbacks
		(record_id, stars, feedback)
		VALUES($1, $2, $3)
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
		WHERE record_id = $1
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
