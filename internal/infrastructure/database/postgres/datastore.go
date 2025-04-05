package postgres

import (
	"context"
	"fmt"
	"timeline/internal/infrastructure/models/datastore"

	"go.uber.org/zap"
)

func (p *PostgresRepo) LoadCities(ctx context.Context, logger *zap.Logger, cities datastore.Cities) error {
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
		INSERT INTO cities (name, tzid)
		VALUES ($1, $2);
	`
	quarter := len(cities.Arr) / 4
	percent := 25
	city, ok := cities.Next()
	for cnt := 1; ok; cnt++ {
		city, ok = cities.Next()
		res, err := tx.ExecContext(ctx, query, city.Name, city.TZ)
		switch {
		case err != nil:
			return fmt.Errorf("failed to save cities: %w", err)
		default:
			if rowsAffected, _ := res.RowsAffected(); rowsAffected == 0 {
				return ErrNoRowsAffected
			}
		}
		if cnt == quarter || cnt == 2*quarter || cnt == 3*quarter || cnt == 4*quarter {
			logger.Info(fmt.Sprintf("Saved %d%% cities", percent))
			percent += 25
		}
	}
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit tx: %w", err)
	}
	return nil
}
