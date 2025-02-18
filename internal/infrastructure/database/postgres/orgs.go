package postgres

import (
	"context"
	"errors"
	"fmt"
	"timeline/internal/infrastructure/models"
	"timeline/internal/infrastructure/models/orgmodel"

	pgx "github.com/jackc/pgx/v5"
)

var (
	ErrOrgExists    = errors.New("org already exists")
	ErrOrgNotFound  = errors.New("org not found")
	ErrOrgsNotFound = errors.New("orgs not found")
)

// Сохраняет в БД инфу об организации. Отдает org_id при успешном сохранении
// Вернет ошибку если такой пользователь уже существует
func (p *PostgresRepo) OrgSave(ctx context.Context, org *orgmodel.OrgRegister) (int, error) {
	tx, err := p.db.Beginx()
	if err != nil {
		return 0, fmt.Errorf("failed to start tx: %w", err)
	}
	// При возникновении ошибки транзакция откатывается по выходу из функции
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	query := `
		INSERT INTO orgs (uuid, email, passwd_hash, name, type, city, address, telephone, lat, long, about)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
        RETURNING org_id;
	`
	orgID := 0

	if err = tx.QueryRowContext(ctx, query,
		org.UUID,
		org.HashCreds.Email,
		org.HashCreds.PasswdHash,
		org.Name,
		org.Type,
		org.City,
		org.Address,
		org.Telephone,
		org.Lat,
		org.Long,
		org.About,
	).Scan(&orgID); err != nil {
		// TODO: пока отдаем фулл ошибку, а вообще нельзя внутренние ошибки отдавать наружу
		return 0, fmt.Errorf("failed to save org: %w", err)
	}
	if tx.Commit() != nil {
		return 0, fmt.Errorf("failed to commit transaction: %w", err)
	}
	return orgID, nil
}

func (p *PostgresRepo) OrgSaveShowcaseImageURL(ctx context.Context, meta *models.ImageMeta) error {
	tx, err := p.db.Beginx()
	if err != nil {
		return fmt.Errorf("failed to start tx: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	query := `INSERT INTO showcase (url, org_id, type)
		VALUES($1, $2, $3);
	`
	if _, err := tx.ExecContext(ctx, query, meta.URL, meta.DomenID, meta.Type); err != nil {
		return fmt.Errorf("failed to save showcase url: %w", err)
	}
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit tx: %w", err)
	}
	return nil
}
func (p *PostgresRepo) OrgUUID(ctx context.Context, orgID int) (string, error) {
	tx, err := p.db.Beginx()
	if err != nil {
		return "", fmt.Errorf("failed to start tx: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	query := `
		SELECT uuid
		FROM orgs
		WHERE org_id = $1;
	`
	var uuid string
	if err := tx.QueryRowContext(ctx, query, orgID).Scan(&uuid); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", ErrOrgNotFound
		}
		return "", fmt.Errorf("failed to get org uuid by id: %w", err)
	}
	if err = tx.Commit(); err != nil {
		return "", fmt.Errorf("failed to commit tx: %w", err)
	}
	return uuid, nil
}
func (p *PostgresRepo) OrgSetUUID(ctx context.Context, orgID int, NewUUID string) error {
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
		UPDATE orgs
		SET
			uuid = $1
		WHERE org_id = $2;
	`
	if rows, err := tx.ExecContext(ctx, query, NewUUID, orgID); err != nil {
		if _, errNoRowsAffected := rows.RowsAffected(); errNoRowsAffected != nil {
			return fmt.Errorf("nothing was added: %w", err)
		}
		return fmt.Errorf("failed to set org uuid: %w", err)
	}
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit tx: %w", err)
	}
	return nil
}

// Если GALLERY/BANNER = true, то строка удаляется, иначе uuid организации будет пустым
func (p *PostgresRepo) OrgDeleteURL(ctx context.Context, meta *models.ImageMeta) error {
	tx, err := p.db.Beginx()
	if err != nil {
		return fmt.Errorf("failed to start tx: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	var query string
	var entity string
	switch {
	case (meta.Type == "gallery") || (meta.Type == "banner"):
		entity = meta.Type
		query = `DELETE FROM showcase WHERE url = $1;`
	case (meta.Type == "org"):
		entity = "org"
		query = `UPDATE orgs SET uuid = '' WHERE uuid = $1;`
	default:
		return fmt.Errorf("image type %s doesn't exist: %w", entity, err)
	}
	if rows, err := tx.ExecContext(ctx, query, meta.URL); err != nil {
		if _, errNoRowsAffected := rows.RowsAffected(); errNoRowsAffected != nil {
			return fmt.Errorf("failed %s urls wasn't deleted: %w", entity, err)
		}
		return fmt.Errorf("failed to delete %s urls: %w", entity, err)
	}
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit tx: %w", err)
	}
	return nil
}

func (p *PostgresRepo) OrgByID(ctx context.Context, id int) (*orgmodel.Organization, error) {
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
		SELECT org_id, email, uuid, name, rating, type, city, address, telephone, lat, long, about 
		FROM orgs
        WHERE is_delete = false 
		AND org_id = $1;
	`
	var org orgmodel.Organization
	if err = tx.GetContext(ctx, &org, query, id); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrOrgNotFound
		}
		return nil, fmt.Errorf("failed to get org by id: %w", err)
	}
	query = `
		SELECT weekday, open, close, break_start, break_end
		FROM timetables
		WHERE org_id = $1;
	`
	if err := tx.SelectContext(ctx, &org.Timetable, query, id); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrOrgNotFound
		}
		return nil, fmt.Errorf("failed to get org timetable by id: %w", err)
	}
	query = `
		SELECT 
			url, type
		FROM showcase
		WHERE org_id = $1;
	`
	if err := tx.SelectContext(ctx, &org.ShowcasesURL, query, id); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrOrgNotFound
		}
		return nil, fmt.Errorf("failed to get org timetable by id: %w", err)
	}
	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit tx: %w", err)
	}
	return &org, nil
}

// Принимает две точки на карте по диагонале. Возвращает список организаций лежащий в диапазоне этих двух точек
func (p *PostgresRepo) OrgsInArea(ctx context.Context, area *orgmodel.AreaParams) ([]*orgmodel.OrgByArea, error) {
	tx, err := p.db.Beginx()
	if err != nil {
		return nil, fmt.Errorf("failed to start tx: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	// Запрос с расписанием на текущий день
	query := `SELECT
		o.org_id,
		o.name,
		o.rating,
		o.type,
		o.lat,
		o.long,
		t.weekday,
		t.open,   
		t.close,
		t.break_start, 
		t.break_end
	FROM orgs o
	LEFT JOIN timetables t
		ON o.org_id = t.org_id
		AND t.weekday = EXTRACT(ISODOW FROM CURRENT_DATE)
	WHERE o.is_delete = false 
	AND o.lat BETWEEN $1 AND $2
	AND o.long BETWEEN $3 AND $4;
	`

	// Результирующий список
	orgList := make([]*orgmodel.OrgByArea, 0, 1)
	if err = tx.SelectContext(
		ctx,
		&orgList,
		query,
		area.Left.Lat,   // Нижняя граница широты
		area.Right.Lat,  // Верхняя граница широты
		area.Left.Long,  // Левая граница долготы
		area.Right.Long, // Правая граница долготы
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrOrgsNotFound
		}
		return nil, fmt.Errorf("failed to get orgs in area with schedule: %w", err)
	}

	// Завершаем транзакцию
	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit tx: %w", err)
	}

	return orgList, nil
}

// Принимает пагинацию, имя и тип организации. Возвращает соответствующие организации, число найденных
func (p *PostgresRepo) OrgsBySearch(ctx context.Context, params *orgmodel.SearchParams) ([]*orgmodel.OrgsBySearch, int, error) {
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
			COUNT(org_id)
		FROM orgs
		WHERE is_delete = false 
		AND ($1 = '' OR name ILIKE '%' || $1 || '%') 
		AND ($2 = '' OR type ILIKE '%' || $2 || '%')
	`
	var found int
	if err = tx.QueryRowxContext(ctx, query, params.Name, params.Type).Scan(&found); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, 0, ErrOrgsNotFound
		}
		return nil, 0, fmt.Errorf("failed orgs searching: %w", err)
	}
	query = `SELECT 
			o.org_id,
			o.name,
			o.rating,
			o.type,
			o.address,
			o.lat,
			o.long,
			t.weekday,
			t.open,   
			t.close,
			t.break_start, 
			t.break_end
		FROM orgs o
		LEFT JOIN timetables t
		ON t.org_id = o.org_id
		AND t.weekday = EXTRACT(ISODOW FROM CURRENT_DATE)
		WHERE o.is_delete = false
		AND ($1 = '' OR name ILIKE '%' || $1 || '%')
		AND ($2 = '' OR type ILIKE '%' || $2 || '%')
		LIMIT $3
		OFFSET $4;
	`
	orgList := make([]*orgmodel.OrgsBySearch, 0, 3)
	if err = tx.SelectContext(ctx, &orgList, query, params.Name, params.Type, params.Limit, params.Offset); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, 0, ErrOrgsNotFound
		}
		return nil, 0, fmt.Errorf("failed orgs searching: %w", err)
	}
	if err = tx.Commit(); err != nil {
		return nil, 0, fmt.Errorf("failed to commit tx: %w", err)
	}
	return orgList, found, nil
}

func (p *PostgresRepo) OrgUpdate(ctx context.Context, new *orgmodel.Organization) error {
	tx, err := p.db.Beginx()
	if err != nil {
		return fmt.Errorf("failed to start tx: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	query := `UPDATE orgs
		SET 
			name = $1,
			type = $2,
			city = $3,
			address = $4,
			telephone = $5,
			lat = $6,
			long = $7,
			about = $8
		WHERE is_delete = false
		AND org_id = $9;
	`
	res, err := tx.ExecContext(ctx, query,
		new.Name,
		new.Type,
		new.City,
		new.Address,
		new.Telephone,
		new.Lat,
		new.Long,
		new.About,
		new.OrgID,
	)
	if err != nil {
		if _, errNoRowsAffected := res.RowsAffected(); errNoRowsAffected != nil {
			return ErrOrgNotFound
		}
		return fmt.Errorf("failed update org info: %w", err)
	}
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit tx: %w", err)
	}
	return nil
}

// [CRON]: удаление стухших юзер аккаунтов
func (p *PostgresRepo) OrgDeleteExpired(ctx context.Context) error {
	tx, err := p.db.Beginx()
	if err != nil {
		return fmt.Errorf("failed to start tx: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	query := `DELETE FROM orgs
		WHERE verified IS NOT true
		AND created_at < CURRENT_TIMESTAMP - INTERVAL '1 day';
	`
	if _, err := tx.ExecContext(ctx, query); err != nil {
		return err
	}
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit tx: %w", err)
	}
	return nil
}
