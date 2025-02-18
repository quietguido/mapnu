package booking

import (
	"context"
	"strconv"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/quietguido/mapnu/internal/repo/booking/model"
	"github.com/quietguido/mapnu/pkg/assert"
)

const (
	bookingTable = "bookings"
)

type repository struct {
	lg      *zap.Logger
	db      *sqlx.DB
	builder sq.StatementBuilderType
}

func NewRepository(lg *zap.Logger, db *sqlx.DB) *repository {
	return &repository{
		lg:      lg,
		db:      db,
		builder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (rp *repository) CreateBooking(ctx context.Context, createBooking model.CreateBooking) (int, error) {
	insertQuery := rp.builder.
		Insert(bookingTable).Columns(
		"user_id",
		"event_id",
		"visibility",
	).Values(
		createBooking.UserID,
		createBooking.EventID,
		createBooking.Visibility,
	).Suffix("RETURNING booking_id")

	sql, args, err := insertQuery.ToSql()
	assert.IsNil(err, "Failed to build SQL query")

	var bookingID int
	err = rp.db.QueryRowContext(ctx, sql, args...).Scan(&bookingID)
	if err != nil {
		rp.lg.Warn(sql)
		return 0, errors.Wrap(err, "Failed to execute SQL query")
	}

	return bookingID, nil
}

func (rp *repository) GetBookingById(ctx context.Context, bookingId int) (*model.Booking, error) {
	selectQuery := `
		SELECT
			booking_id,
			user_id,
			event_id,
			booking_status,
			visibility,
			booked_at
		FROM bookings
		WHERE booking_id = $1;
	`

	// Execute the query
	row := rp.db.QueryRowxContext(ctx, selectQuery, bookingId)
	var booking model.Booking
	err := row.StructScan(&booking)
	if err != nil {
		rp.lg.Error("SQL Query Failed:", zap.String("query", selectQuery))
		rp.lg.Error("Booking ID:", zap.String("booking_id", strconv.Itoa(bookingId)))
		return nil, errors.Wrap(err, "Failed to execute SQL query")
	}
	return &booking, nil
}

func (rp *repository) GetBookingsForUser(ctx context.Context, userID uuid.UUID) ([]model.Booking, error) {
	selectQuery := `
		SELECT
			booking_id,
			user_id,
			event_id,
			booking_status,
			visibility,
			booked_at
		FROM bookings
		WHERE user_id = $1
		ORDER BY booked_at DESC; -- Sort by latest booking first
	`

	// Execute the query
	rows, err := rp.db.QueryxContext(ctx, selectQuery, userID)
	if err != nil {
		rp.lg.Error("Failed to execute GetBookingsForUser query", zap.Error(err))
		return nil, errors.Wrap(err, "Failed to execute SQL query")
	}
	defer rows.Close()

	var bookings []model.Booking
	for rows.Next() {
		var booking model.Booking
		err := rows.StructScan(&booking)
		if err != nil {
			rp.lg.Error("Failed to scan row", zap.Error(err))
			return nil, errors.Wrap(err, "Failed to scan row")
		}
		bookings = append(bookings, booking)
	}

	if err := rows.Err(); err != nil {
		rp.lg.Error("Row iteration error", zap.Error(err))
		return nil, errors.Wrap(err, "Row iteration error")
	}

	return bookings, nil
}

func (rp *repository) ChangeBookingStatus(ctx context.Context, bookingId int, status string) error {
	updateQuery := rp.builder.
		Update(bookingTable).
		Set("booking_status", status).
		Where(sq.Eq{"booking_id": bookingId})

	sql, args, err := updateQuery.ToSql()
	if err != nil {
		rp.lg.Error("Failed to build SQL query", zap.Error(err))
		return errors.Wrap(err, "Failed to build SQL query")
	}

	result, err := rp.db.ExecContext(ctx, sql, args...)
	if err != nil {
		rp.lg.Error("Failed to execute SQL query", zap.Error(err))
		return errors.Wrap(err, "Failed to update status")
	}

	num, err := result.RowsAffected()
	if err != nil {
		rp.lg.Error("Failed to get affected rows", zap.Error(err))
		return errors.Wrap(err, "Failed to update status")
	}
	if num == 0 {
		rp.lg.Warn("No booking found with the given ID", zap.Int("booking_id", bookingId))
		return errors.New("No booking found with the given ID")
	}

	return nil
}
