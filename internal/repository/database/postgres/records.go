package postgres

import (
	"context"
	"fmt"
	"timeline/internal/repository/models/orgmodel"
	"timeline/internal/repository/models/recordmodel"
	"timeline/internal/repository/models/usermodel"
)

func (p *PostgresRepo) Record(ctx context.Context, recordID int) (*recordmodel.RecordScrap, error) {
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
		SELECT 
		srvc.name AS service_name, 
		srvc.cost, 
		w.first_name AS worker_first_name, 
		w.last_name AS worker_last_name, 
		o.name AS org_name,
		u.first_name AS user_first_name,
		u.last_name AS user_last_name, 
		s.date,
		s.session_begin,
		s.session_end,
		f.stars,
		f.feedback,
		r.reviewed
	FROM records r
	JOIN slots s ON r.slot_id = s.slot_id
	JOIN orgs o ON r.org_id = o.org_id
	JOIN users u ON r.user_id = u.user_id
	JOIN services srvc ON r.service_id = srvc.service_id
	JOIN workers w ON r.worker_id = w.worker_id
	LEFT JOIN feedbacks f ON r.record_id = f.record_id
	WHERE r.record_id = $1;
	`
	rec := &recordmodel.RecordScrap{
		RecordID: recordID,
		Org:      &orgmodel.OrgInfo{},
		User:     &usermodel.UserInfo{},
		Slot:     &orgmodel.Slot{},
		Service:  &orgmodel.Service{},
		Worker:   &orgmodel.Worker{},
		Feedback: &recordmodel.Feedback{},
	}
	if err := tx.QueryRowxContext(ctx, query, recordID).Scan(
		&rec.Service.Name,
		&rec.Service.Cost,
		&rec.Worker.FirstName,
		&rec.Worker.LastName,
		&rec.Org.Name,
		&rec.User.FirstName,
		&rec.User.LastName,
		&rec.Slot.Date,
		&rec.Slot.Begin,
		&rec.Slot.End,
		&rec.Feedback.Stars,
		&rec.Feedback.Feedback,
		&rec.Reviewed,
	); err != nil {
		return nil, err
	}
	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit tx: %w", err)
	}
	return rec, nil
}

func (p *PostgresRepo) RecordList(ctx context.Context, req *recordmodel.RecordListParams) ([]*recordmodel.RecordScrap, error) {
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
		SELECT 
		srvc.name AS service_name, 
		srvc.cost, 
		w.first_name AS worker_first_name, 
		w.last_name AS worker_last_name, 
		o.name AS org_name,
		u.first_name AS user_first_name,
		u.last_name AS user_last_name, 
		s.date,
		s.session_begin,
		s.session_end,
		f.stars,
		f.feedback,
		r.reviewed,
		r.record_id
	FROM records r
	JOIN slots s ON r.slot_id = s.slot_id
	JOIN orgs o ON r.org_id = o.org_id
	JOIN users u ON r.user_id = u.user_id
	JOIN services srvc ON r.service_id = srvc.service_id
	JOIN workers w ON r.worker_id = w.worker_id
	LEFT JOIN feedbacks f ON r.record_id = f.record_id
	WHERE r.user_id = $1
	OR r.org_id = $2;
	`
	recs := make([]*recordmodel.RecordScrap, 0, 2)
	rows, err := tx.QueryContext(ctx, query, req.UserID, req.OrgID)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	rec := &recordmodel.RecordScrap{
		Org:      &orgmodel.OrgInfo{},
		User:     &usermodel.UserInfo{},
		Slot:     &orgmodel.Slot{},
		Service:  &orgmodel.Service{},
		Worker:   &orgmodel.Worker{},
		Feedback: &recordmodel.Feedback{},
	}
	for rows.Next() {
		err := rows.Scan(
			&rec.Service.Name,
			&rec.Service.Cost,
			&rec.Worker.FirstName,
			&rec.Worker.LastName,
			&rec.Org.Name,
			&rec.User.FirstName,
			&rec.User.LastName,
			&rec.Slot.Date,
			&rec.Slot.Begin,
			&rec.Slot.End,
			&rec.Feedback.Stars,
			&rec.Feedback.Feedback,
			&rec.Reviewed,
			&rec.RecordID,
		)
		if err != nil {
			return nil, err
		}
		recs = append(recs, rec)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit tx: %w", err)
	}
	return recs, nil
}

func (p *PostgresRepo) RecordAdd(ctx context.Context, req *recordmodel.Record) error {
	tx, err := p.db.Beginx()
	if err != nil {
		return fmt.Errorf("failed to start tx: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	query := `INSERT INTO records
		(record_id, user_id, org_id, service_id, slot_id, worker_id, reviewed)
		VALUES($1, $2, $3, $4, $5, $6, $7)
	`
	rows, err := tx.ExecContext(
		ctx,
		query,
		req.RecordID,
		req.UserID,
		req.OrgID,
		req.ServiceID,
		req.SlotID,
		req.WorkerID,
		req.Reviewed,
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

func (p *PostgresRepo) RecordPatch(ctx context.Context, req *recordmodel.Record) error {
	tx, err := p.db.Beginx()
	if err != nil {
		return fmt.Errorf("failed to start tx: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	query := `UPDATE records
		SET
			user_id = COALESCE(NULLIF($1, 0), user_id),
			org_id = COALESCE(NULLIF($2, 0), org_id),
			service_id = COALESCE(NULLIF($3, 0), service_id),
			slot_id = COALESCE(NULLIF($4, 0), slot_id),
			worker_id = COALESCE(NULLIF($5, 0), worker_id),
			reviewed = COALESCE($6, reviewed)
		WHERE record_id = $7
	`
	rows, err := tx.ExecContext(
		ctx,
		query,
		req.UserID,
		req.OrgID,
		req.ServiceID,
		req.SlotID,
		req.WorkerID,
		req.Reviewed,
		req.RecordID,
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

func (p *PostgresRepo) RecordDelete(ctx context.Context, req *recordmodel.Record) error {
	tx, err := p.db.Beginx()
	if err != nil {
		return fmt.Errorf("failed to start tx: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	query := `DELETE FROM records
		WHERE record_id 
		IN (SELECT r.record_id
			FROM records r
			JOIN slots s ON r.slot_id = s.slot_id AND r.record_id = $1
			WHERE s.date >= CURRENT_DATE
			AND CURRENT_TIME < (s.session_begin::time - INTERVAL '2 hours')
		);
	`
	rows, err := tx.ExecContext(
		ctx,
		query,
		req.RecordID,
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
