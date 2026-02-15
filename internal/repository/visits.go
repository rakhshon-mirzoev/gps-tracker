package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/model"
	"github.com/rakhshon-mirzoev/gps-tracker/pkg/db"
	errs "github.com/rakhshon-mirzoev/gps-tracker/pkg/error"
)

type VisitRepository interface {
	Create(ctx context.Context, v model.Visit) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.Visit, error)
	List(ctx context.Context, scheduleID uuid.UUID) ([]model.Visit, error)
	ListByStaffID(ctx context.Context, staff_id uuid.UUID, lim, off int) ([]model.Visit, int, int, error)
	ListByReviewStatus(ctx context.Context, status string, limit, offset int, id uuid.UUID) ([]model.VisitResponse, int, error)
	Update(ctx context.Context, id uuid.UUID, dto model.UpdateVisit) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetDataForExcel(ctx context.Context, staff_id, store_id *uuid.UUID, begin, end *time.Time) ([]model.VisitExcel, error)
	GetByStoreScheduleId(ctx context.Context, storeID, scheduleID, staffID uuid.UUID) (*model.Visit, error)
	CreateTx(ctx context.Context, tx pgx.Tx, v *model.Visit) error
	BeginTx(ctx context.Context) (pgx.Tx, error)
}

type visit struct {
	DB *db.DB
}

func NewVisitRepository(db *db.DB) VisitRepository {
	return &visit{DB: db}
}

func (r *visit) BeginTx(ctx context.Context) (pgx.Tx, error) {
	return r.DB.Pool.BeginTx(ctx, pgx.TxOptions{
		IsoLevel: pgx.ReadCommitted,
	})

}

func (r *visit) CreateTx(ctx context.Context, tx pgx.Tx, v *model.Visit) error {
	sql := `
		INSERT INTO visits (
			schedule_id,
			staff_id,
			store_id,
			check_in_at,
			check_in_lat,
			check_in_lon,
			distance_deviation,
			is_offline,
			reviewer_id
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9
		)
		RETURNING id, created_at
	`

	err := tx.QueryRow(
		ctx,
		sql,
		v.ScheduleID,
		v.StaffID,
		v.StoreID,
		v.CheckInAt,
		v.CheckInLat,
		v.CheckInLon,
		v.DistanceDeviation,
		v.IsOffline,
		v.ReviewerID,
	).Scan(&v.ID, &v.CreatedAt)

	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			if pgErr.Code == "23505" {
				return errs.ErrroActiveVisitExists
			}
		}
		return fmt.Errorf("create visit tx: %w", err)
	}

	return nil
}

func (r *visit) Create(ctx context.Context, v model.Visit) error {
	sql := `
        INSERT INTO visits (
            schedule_id,
            staff_id,
            store_id,
            check_in_at,
            check_in_lat,
			check_in_lon
        ) VALUES (
            $1, $2, $3, $4, $5, $6
    	);
    `

	err := r.DB.Pool.QueryRow(ctx, sql,
		v.ScheduleID, v.StaffID, v.StoreID, v.CheckInAt,
		v.CheckInLat, v.CheckInLon,
	).Scan(&v.ID, &v.CreatedAt)
	if err != nil {
		return fmt.Errorf("create visit: %w", err)
	}

	return nil
}

func (r *visit) GetByID(ctx context.Context, id uuid.UUID) (*model.Visit, error) {
	sql := `
        SELECT id, schedule_id, staff_id, store_id, check_in_at, check_out_at, 
			   check_in_lat, check_in_lon, distance_deviation, comment,
               is_offline, review_status, review_comment, reviewer_id,
               reviewed_at, created_at
        FROM visits WHERE id = $1
    `

	var v model.Visit
	err := r.DB.Pool.QueryRow(ctx, sql, id).Scan(
		&v.ID, &v.ScheduleID, &v.StaffID, &v.StoreID, &v.CheckInAt, &v.CheckOutAt,
		&v.CheckInLat, &v.CheckInLon, &v.DistanceDeviation, &v.Comment,
		&v.IsOffline, &v.ReviewStatus, &v.ReviewComment, &v.ReviewerID,
		&v.ReviewedAt, &v.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("get by id visit: %w", err)
	}

	return &v, nil
}

func (r *visit) List(ctx context.Context, scheduleID uuid.UUID) ([]model.Visit, error) {
	sql := `
        SELECT id, schedule_id, staff_id, store_id, check_in_at, check_out_at,
               check_in_lat, check_in_lon, distance_deviation, comment,
               is_offline, review_status, review_comment, reviewer_id,
               reviewed_at, created_at
        FROM visits WHERE schedule_id = $1
    `

	rows, err := r.DB.Pool.Query(ctx, sql, scheduleID)
	if err != nil {
		return nil, fmt.Errorf("list visits: %w", err)
	}

	var visits = make([]model.Visit, 0)

	for rows.Next() {
		var v model.Visit

		if err := rows.Scan(
			&v.ID, &v.ScheduleID, &v.StaffID, &v.StoreID, &v.CheckInAt, &v.CheckOutAt,
			&v.CheckInLat, &v.CheckInLon, &v.DistanceDeviation, &v.Comment,
			&v.IsOffline, &v.ReviewStatus, &v.ReviewComment, &v.ReviewerID,
			&v.ReviewedAt, &v.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan visit: %w", err)
		}

		visits = append(visits, v)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows visits: %w", err)
	}

	return visits, nil
}

func (r *visit) Update(ctx context.Context, id uuid.UUID, dto model.UpdateVisit) error {
	sql := `UPDATE visits SET `
	var args []any
	n := 1

	if dto.ScheduleID != nil {
		sql += fmt.Sprintf("schedule_id = $%d, ", n)
		args = append(args, *dto.ScheduleID)
		n++
	}

	if dto.StoreID != nil {
		sql += fmt.Sprintf("store_id = $%d, ", n)
		args = append(args, *dto.StoreID)
		n++
	}

	if dto.CheckInAt != nil {
		sql += fmt.Sprintf("check_in_at = $%d, ", n)
		args = append(args, *dto.CheckInAt)
		n++
	}

	if dto.CheckOutAt != nil {
		sql += fmt.Sprintf("check_out_at = $%d, ", n)
		args = append(args, *dto.CheckOutAt)
		n++
	}

	if dto.CheckInLat != nil {
		sql += fmt.Sprintf("check_in_lat = $%d, ", n)
		args = append(args, *dto.CheckInLat)
		n++
	}

	if dto.CheckInLon != nil {
		sql += fmt.Sprintf("check_in_lon = $%d, ", n)
		args = append(args, *dto.CheckInLon)
		n++
	}

	if dto.DistanceDeviation != nil {
		sql += fmt.Sprintf("distance_deviation = $%d, ", n)
		args = append(args, *dto.DistanceDeviation)
		n++
	}

	if dto.Comment != nil {
		sql += fmt.Sprintf("comment = $%d, ", n)
		args = append(args, *dto.Comment)
		n++
	}

	if dto.IsOffline != nil {
		sql += fmt.Sprintf("is_offline = $%d, ", n)
		args = append(args, *dto.IsOffline)
		n++
	}

	if dto.ReviewStatus != nil {
		sql += fmt.Sprintf("review_status = $%d, ", n)
		args = append(args, *dto.ReviewStatus)
		n++
	}

	if dto.ReviewComment != nil {
		sql += fmt.Sprintf("review_comment = $%d, ", n)
		args = append(args, *dto.ReviewComment)
		n++
	}

	if dto.ReviewerID != nil {
		sql += fmt.Sprintf("reviewer_id = $%d, ", n)
		args = append(args, *dto.ReviewerID)
		n++
	}

	if dto.ReviewedAt != nil {
		sql += fmt.Sprintf("reviewed_at = $%d, ", n)
		args = append(args, *dto.ReviewedAt)
		n++
	}

	if len(args) == 0 {
		return nil
	}

	sql = strings.TrimSuffix(sql, ", ")
	sql += fmt.Sprintf(" WHERE id = $%d", n)
	args = append(args, id)

	_, err := r.DB.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("update visit: %w", err)
	}

	return nil
}

func (r *visit) ListByReviewStatus(ctx context.Context, status string, limit, offset int, id uuid.UUID) ([]model.VisitResponse, int, error) {
	const defaultLimit = 20

	if limit <= 0 {
		limit = defaultLimit
	}
	if offset < 0 {
		offset = 0
	}

	var total int
	countSQL := `SELECT COUNT(*) FROM visits WHERE review_status = $1 AND reviewer_id = $2`
	if err := r.DB.Pool.QueryRow(ctx, countSQL, status, id).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count visits: %w", err)
	}

	sql := `
		SELECT v.id, v.schedule_id, v.staff_id, v.store_id, v.check_in_at, v.check_out_at,
			   v.check_in_lat, v.check_in_lon, v.distance_deviation, v.comment,
			   v.is_offline, v.review_status, v.review_comment, v.reviewer_id,
			   v.reviewed_at, v.created_at,
			   su.id, su.full_name, su.phone_number,
			   st.id, st.name, st.address, st.latitude, st.longitude
		FROM visits v
		LEFT JOIN staff_users su ON v.staff_id = su.id
		LEFT JOIN stores st ON v.store_id = st.id
		WHERE v.review_status = $1 AND v.reviewer_id = $2
		ORDER BY v.created_at DESC
		LIMIT $3 OFFSET $4
	`

	rows, err := r.DB.Pool.Query(ctx, sql, status, id, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("list visits by review status: %w", err)
	}
	defer rows.Close()

	var visits = make([]model.VisitResponse, 0)
	for rows.Next() {
		var v model.VisitResponse
		var staff model.VisitStaffBrief
		var store model.VisitStoreBrief
		if err := rows.Scan(
			&v.ID, &v.ScheduleID, &v.StaffID, &v.StoreID, &v.CheckInAt, &v.CheckOutAt,
			&v.CheckInLat, &v.CheckInLon, &v.DistanceDeviation, &v.Comment,
			&v.IsOffline, &v.ReviewStatus, &v.ReviewComment, &v.ReviewerID,
			&v.ReviewedAt, &v.CreatedAt,
			&staff.ID, &staff.FullName, &staff.PhoneNumber,
			&store.ID, &store.Name, &store.Address, &store.Latitude, &store.Longitude,
		); err != nil {
			return nil, 0, fmt.Errorf("scan visit review status: %w", err)
		}
		v.Staff = &staff
		v.Store = &store
		v.ShelfPhotos = make([]string, 0)
		visits = append(visits, v)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("rows visits review status: %w", err)
	}

	return visits, total, nil
}

func (r *visit) Delete(ctx context.Context, id uuid.UUID) error {
	sql := `
        DELETE FROM visits WHERE id = $1
    `

	_, err := r.DB.Pool.Exec(ctx, sql, id)
	if err != nil {
		return fmt.Errorf("delete visit: %w", err)
	}

	return nil
}

func (s *visit) GetDataForExcel(ctx context.Context, staff_id, store_id *uuid.UUID, begin, end *time.Time) ([]model.VisitExcel, error) {
	var conditions []string
	var args []any
	argNum := 1

	if staff_id != nil {
		conditions = append(conditions, fmt.Sprintf("v.staff_id = $%d", argNum))
		args = append(args, *staff_id)
		argNum++
	}

	if store_id != nil {
		conditions = append(conditions, fmt.Sprintf("v.store_id = $%d", argNum))
		args = append(args, *store_id)
		argNum++
	}

	if begin != nil {
		conditions = append(conditions, fmt.Sprintf("v.check_in_at >= $%d", argNum))
		args = append(args, *begin)
		argNum++
	}

	if end != nil {
		conditions = append(conditions, fmt.Sprintf("v.check_in_at < $%d", argNum))
		args = append(args, *end)
		argNum++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	sql := fmt.Sprintf(`
	SELECT
		v.id,
		v.check_in_at,
		v.check_out_at,
		v.check_in_lat,
		v.check_in_lon,
		v.distance_deviation,
		v.comment,
		v.review_status,
		v.review_comment,
		v.reviewed_at,
		v.created_at,

		s.full_name,
		s.phone_number,
		g.name,
		st.name,
		su.full_name

	FROM visits v
	JOIN staff_users s ON v.staff_id = s.id
	JOIN groups g ON g.id = s.group_id
	JOIN stores st ON st.id = v.store_id
	LEFT JOIN system_users su ON su.id = v.reviewer_id
	%s
	ORDER BY v.check_in_at DESC
`, whereClause)

	rows, err := s.DB.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("get data for excel: %w", err)
	}
	defer rows.Close()

	var visits []model.VisitExcel

	for rows.Next() {
		var v model.VisitExcel
		if err := rows.Scan(
			&v.ID,
			&v.CheckInAt,
			&v.CheckOutAt,
			&v.CheckInLat,
			&v.CheckInLon,
			&v.DistanceDeviation,
			&v.Comment,
			&v.ReviewStatus,
			&v.ReviewComment,
			&v.ReviewedAt,
			&v.CreatedAt,
			&v.Staff.FullName,
			&v.Staff.PhoneNumber,
			&v.Group.Name,
			&v.Store.Name,
			&v.Reviewer.FullName,
		); err != nil {
			return nil, fmt.Errorf("scan visit excel: %w", err)
		}

		visits = append(visits, v)

	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows visits excel: %w", err)
	}

	return visits, nil
}

func (s *visit) GetByStoreScheduleId(ctx context.Context, storeID, scheduleID, staffID uuid.UUID) (*model.Visit, error) {
	sql := `
		SELECT id, check_in_at, check_out_at,
			   check_in_lat, check_in_lon, distance_deviation, comment,
			   is_offline, review_status, review_comment, reviewer_id,
			   reviewed_at, created_at
		FROM visits 
		WHERE store_id = $1	AND schedule_id = $2 AND staff_id = $3
	`

	var v model.Visit

	if err := s.DB.Pool.QueryRow(ctx, sql, storeID, scheduleID, staffID).Scan(
		&v.ID,
		&v.CheckInAt,
		&v.CheckOutAt,
		&v.CheckInLat,
		&v.CheckInLon,
		&v.DistanceDeviation,
		&v.Comment,
		&v.IsOffline,
		&v.ReviewStatus,
		&v.ReviewComment,
		&v.ReviewerID,
		&v.ReviewedAt,
		&v.CreatedAt,
	); err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("scan get by store id visit: %w", err)
	}

	return &v, nil
}

func (s *visit) ListByStaffID(ctx context.Context, staff_id uuid.UUID, lim, off int) ([]model.Visit, int, int, error) {
	var (
		v     []model.Visit
		total int
		done  int
	)

	totalSql := `SELECT COUNT(*) FROM visits WHERE staff_id = $1`

	if err := s.DB.Pool.QueryRow(ctx, totalSql, staff_id).Scan(&total); err != nil {
		return nil, 0, 0, fmt.Errorf("count visits: %w", err)
	}

	doneSql := `SELECT COUNT(*) FROM visits where staff_id = $1 AND review_status = 'approved'`

	if err := s.DB.Pool.QueryRow(ctx, doneSql, staff_id).Scan(&done); err != nil {
		return nil, 0, 0, fmt.Errorf("count done: %w", err)
	}

	sql := `
		SELECT id, check_in_at, check_out_at,
			   check_in_lat, check_in_lon, distance_deviation, comment,
			   is_offline, review_status, review_comment, reviewer_id,
			   reviewed_at, created_at
		FROM visits 
		WHERE staff_id = $1
	`

	rows, err := s.DB.Pool.Query(ctx, sql, staff_id)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("query exec list by staff id: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var visit model.Visit

		if err := rows.Scan(
			&visit.ID,
			&visit.CheckInAt,
			&visit.CheckOutAt,
			&visit.CheckInLat,
			&visit.CheckInLon,
			&visit.DistanceDeviation,
			&visit.Comment,
			&visit.IsOffline,
			&visit.ReviewStatus,
			&visit.ReviewComment,
			&visit.ReviewerID,
			&visit.ReviewedAt,
			&visit.CreatedAt,
		); err != nil {
			return nil, 0, 0, fmt.Errorf("rows scan list by staff id: %w", err)
		}

		v = append(v, visit)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, 0, fmt.Errorf("rows err: %w", err)
	}

	return v, total, done, nil
}
