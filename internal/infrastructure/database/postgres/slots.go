package postgres

import (
	"context"
	"fmt"
	"time"
	"timeline/internal/infrastructure/models/orgmodel"
	"timeline/internal/libs/custom"
)

// [CRON]:
// Генерирует свободные слоты на след неделю относительно текущего дня.
func (p *PostgresRepo) GenerateSlots(ctx context.Context) error {
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
		SELECT ws.worker_schedule_id, ws.worker_id, ws.org_id, ws.weekday, ws.start, ws.over, w.session_duration, t.break_start, t.break_end
		FROM worker_schedules ws
		JOIN workers w ON w.worker_id = ws.worker_id
		JOIN timetables t ON t.org_id = ws.org_id AND ws.weekday = t.weekday
		WHERE w.is_delete = false
		AND ws.is_delete = false;
	`
	workerSchedules := make([]*struct {
		WorkerScheduleID int       `db:"worker_schedule_id"`
		WorkerID         int       `db:"worker_id"`
		OrgID            int       `db:"org_id"`
		Weekday          int       `db:"weekday"`
		Start            time.Time `db:"start"`
		Over             time.Time `db:"over"`
		SessionDuration  int       `db:"session_duration"`
		BreakStart       time.Time `db:"break_start"`
		BreakEnd         time.Time `db:"break_end"`
	}, 0, 10)
	if err = tx.SelectContext(ctx, &workerSchedules, query); err != nil {
		return err
	}
	// формула для вычисления будущей даты:
	// new_date = curr_date + should_be_added
	// should_be_added = (7-curr_weekday) + given_weekday + (curr_weekday-1)*7
	// e.g: curr_weekday = 1 (Monday), given_weekday = 3 (Wednesday) =>
	// 		=> should_be_add = (7-1)+3+(1-1)*7 = 9 => 1 + 9 = 10.
	//		день на след неделе: 10%7 = 3 - (Wednesday) - OK
	query = `
		INSERT INTO slots
		(date, session_begin, session_end, busy, worker_schedule_id, worker_id)
		VALUES
		(CURRENT_DATE + ((7 - EXTRACT(ISODOW FROM CURRENT_DATE)) + $1  + ((EXTRACT(ISODOW FROM CURRENT_DATE)-1) * 7) || ' days' )::INTERVAL, 
		$2, 
		$3, 
		$4, 
		$5, 
		$6);
	`
	busy := false
	for _, v := range workerSchedules {
		periods := int(v.Over.Sub(v.Start) / (time.Duration(v.SessionDuration) * time.Minute))
		for i := range periods {
			begin := v.Start.Add(time.Duration(i) * time.Duration(v.SessionDuration) * time.Minute) // begin := start + i*session_duration. e.g: 12:00 + 0*60 = 12:00
			end := v.Start.Add(time.Duration(i+1) * time.Duration(v.SessionDuration) * time.Minute) // end := start + (i+1)*session_duration. e.g: 12:00 + 1*60 = 13:00
			// Слот не создается, если время начала или конца сеанса попадает в перерыв работника
			if custom.CompareTime(begin, v.BreakStart) >= 0 && custom.CompareTime(begin, v.BreakEnd) < 0 {
				continue
			}
			if custom.CompareTime(end, v.BreakStart) > 0 && custom.CompareTime(end, v.BreakEnd) <= 0 {
				continue
			}
			res, err := tx.ExecContext(ctx, query, v.Weekday, begin.UTC(), end.UTC(), busy, v.WorkerScheduleID, v.WorkerID) //nolint:govet // ...
			switch {
			case err != nil:
				return fmt.Errorf("failed to generate slot: %w", err)
			default:
				if rowsAffected, _ := res.RowsAffected(); rowsAffected == 0 {
					return ErrNoRowsAffected
				}
			}
		}
	}
	if tx.Commit() != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}

// [CRON]:
// Удаляет все устаревшие свободные слоты.
// Устаревшими считаются все, что ранее текущего дня
func (p *PostgresRepo) DeleteExpiredSlots(ctx context.Context) error {
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
		DELETE FROM slots
		WHERE date <= CURRENT_DATE
		AND busy = false;
	`
	res, err := tx.ExecContext(ctx, query)
	switch {
	case err != nil:
		return fmt.Errorf("failed to delete expired slots: %w", err)
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

// Занять или освободить слот
func (p *PostgresRepo) UpdateSlot(ctx context.Context, busy bool, params *orgmodel.SlotsMeta) error {
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
		UPDATE slots
		SET
			busy = $1
		WHERE slot_id = $2
		AND worker_id = $3;
	`
	res, err := tx.ExecContext(ctx, query, busy, params.SlotID, params.WorkerID)
	switch {
	case err != nil:
		return fmt.Errorf("failed to take a slot: %w", err)
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

// Получение слотов для указанного работника или его расписания начиная от текущего дня.
func (p *PostgresRepo) Slots(ctx context.Context, params *orgmodel.SlotsMeta) ([]*orgmodel.Slot, error) {
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
		SELECT slot_id, worker_schedule_id, worker_id, date, session_begin, session_end, busy
		FROM slots
		WHERE date >= CURRENT_DATE
		AND ($1 <= 0 OR worker_id = $1);
	`
	slots := make([]*orgmodel.Slot, 0, 1)
	if err = tx.SelectContext(ctx, &slots, query, params.WorkerID); err != nil {
		return nil, err
	}
	if tx.Commit() != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}
	return slots, nil
}
