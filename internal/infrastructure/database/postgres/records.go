package postgres

import (
	"context"
	"errors"
	"fmt"
	"log"
	"timeline/internal/infrastructure/models/orgmodel"
	"timeline/internal/infrastructure/models/recordmodel"
	"timeline/internal/infrastructure/models/usermodel"

	"github.com/jackc/pgx/v5"
)

var (
	ErrRecordsNotFound = errors.New("records not found")
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

func (p *PostgresRepo) RecordList(ctx context.Context, req *recordmodel.RecordListParams) ([]*recordmodel.RecordScrap, int, error) {
	tx, err := p.db.Beginx()
	if err != nil {
		return nil, 0, fmt.Errorf("failed to start tx: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	query := `
		SELECT
			COUNT(*)
		FROM records r
		JOIN slots s ON r.slot_id = s.slot_id
		JOIN orgs o ON r.org_id = o.org_id
		JOIN users u ON r.user_id = u.user_id
		JOIN services srvc ON r.service_id = srvc.service_id
		JOIN workers w ON r.worker_id = w.worker_id
		LEFT JOIN feedbacks f ON r.record_id = f.record_id
		WHERE ($1 <= 0 OR r.user_id = $1)
		AND ($2 <= 0 OR r.org_id = $2)
		AND ( 
				($3 = TRUE AND s.date >= CURRENT_DATE) 
				OR 
				($3 = FALSE AND s.date < CURRENT_DATE)   
			);
	`
	var found int
	if err = tx.QueryRowxContext(ctx, query, req.UserID, req.OrgID, req.Fresh).Scan(&found); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, 0, ErrRecordsNotFound
		}
		return nil, 0, fmt.Errorf("failed to retrieve found: %w", err)
	}
	query = `
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
		WHERE ($1 <= 0 OR r.user_id = $1)
		AND ($2 <= 0 OR r.org_id = $2)
		AND ( 
				($3 = TRUE AND s.date >= CURRENT_DATE) 
				OR 
				($3 = FALSE AND s.date < CURRENT_DATE)   
			)
		LIMIT $4
		OFFSET $5;
	`
	recs := make([]*recordmodel.RecordScrap, 0, 3)
	rows, err := tx.QueryContext(ctx, query, req.UserID, req.OrgID, req.Fresh, req.Limit, req.Offset)
	if err != nil {
		return nil, 0, err
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
		fmt.Println("Entered")
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
			return nil, 0, err
		}
		recs = append(recs, rec)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}
	if err = tx.Commit(); err != nil {
		return nil, 0, fmt.Errorf("failed to commit tx: %w", err)
	}
	return recs, found, nil
}

func (p *PostgresRepo) RecordAdd(ctx context.Context, req *recordmodel.Record) (*recordmodel.ReminderRecord, error) {
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
		INSERT INTO records
			(user_id, org_id, service_id, slot_id, worker_id)
		SELECT $1, $2, $3, $4, $5
		FROM slots s
		WHERE s.slot_id = $4
		AND s.busy = FALSE
		RETURNING record_id;
	`
	var recordID int
	log.Println(recordID)
	if err = tx.QueryRowContext(
		ctx,
		query,
		req.UserID,
		req.OrgID,
		req.ServiceID,
		req.SlotID,
		req.WorkerID,
	).Scan(&recordID); err != nil {
		return nil, err
	}
	query = `
		SELECT 
			u.email AS user_email,
			srvc.name AS service_name,
			srvc.description AS service_description,
			o.name AS org_name,
			o.address AS org_address,
			s.date,
			s.session_begin,
			s.session_end
		FROM records r
		JOIN slots s ON r.slot_id = s.slot_id
		JOIN orgs o ON r.org_id = o.org_id
		JOIN users u ON r.user_id = u.user_id
		JOIN services srvc ON r.service_id = srvc.service_id
		WHERE r.record_id = $1;
	`
	record := &recordmodel.ReminderRecord{}
	if err := tx.QueryRowContext(ctx, query, recordID).Scan(
		&record.UserEmail,
		&record.ServiceName,
		&record.ServiceDescription,
		&record.OrgName,
		&record.OrgAddress,
		&record.Date,
		&record.Begin,
		&record.End,
	); err != nil {
		return nil, err
	}
	fmt.Println(record)
	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit tx: %w", err)
	}
	return record, nil
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
			return fmt.Errorf("no rows inserted") // TODO: везде вот такую строку лучше отдавать, т.к err.Error() бывает пустая в таком случае
		}
	}
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit tx: %w", err)
	}
	return nil
}

func (p *PostgresRepo) UpcomingRecords(ctx context.Context) ([]*recordmodel.ReminderRecord, error) {
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
			u.email AS user_email,
			srvc.name AS service_name,
			srvc.description AS service_description
			o.name AS org_name,
			o.address AS org_address
			s.date,
			s.session_begin,
			s.session_end
		FROM records r
		JOIN slots s ON r.slot_id = s.slot_id
		JOIN orgs o ON r.org_id = o.org_id
		JOIN users u ON r.user_id = u.user_id
		JOIN services srvc ON r.service_id = srvc.service_id
		WHERE s.date = CURRENT_DATE
		AND CURRENT_TIME >= (session_begin::time - INTERVAL '2 hour')
		AND CURRENT_TIME < session_begin::time;
	`
	recs := make([]*recordmodel.ReminderRecord, 0, 2)
	rows, err := tx.QueryContext(ctx, query)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	rec := &recordmodel.ReminderRecord{}
	for rows.Next() {
		err := rows.Scan(
			rec.UserEmail,
			rec.ServiceName,
			rec.ServiceDescription,
			rec.OrgName,
			rec.OrgAddress,
			rec.Date,
			rec.Begin,
			rec.End,
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