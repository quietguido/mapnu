package event

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/quietguido/mapnu/internal/repo/event/model"
)

const (
	eventTable = "event"
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

func (rp *repository) CreateEvent(ctx context.Context, createEvent model.CreateEvent) error {
	// Use Squirrel to build the query
	insertQuery := rp.builder.
		Insert(eventTable).Columns(
		"name",
		"description",
		"created_by",
		"location_date", // PostGIS point with measurement (POINTM)
		"organizer",
		"upvote",
		"downvote",
		// "created_at",
	).Values(
		createEvent.Name,
		createEvent.Description,
		createEvent.CreatedBy,
		createPointM(
			createEvent.Location_lon,
			createEvent.Location_lat,
			createEvent.Time,
		), // Custom function to generate POINTM
		createEvent.Organizer,
		createEvent.Upvote,
		createEvent.Downvote,
		// createEvent.CreatedAt,
	)

	sql, args, err := insertQuery.ToSql()
	if err != nil {
		return errors.Wrap(err, "Failed to build SQL query")
	}

	// Execute the query
	result, err := rp.db.ExecContext(ctx, sql, args...)
	if err != nil {
		return errors.Wrap(err, "Failed to execute SQL query")
	}

	// Optionally check the number of rows affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "Failed to get affected rows")
	}
	if rowsAffected == 0 {
		return errors.New("No rows were inserted")
	}

	return nil
}

func (rp *repository) GetEventById(ctx context.Context, eventId string) (*model.Event, error) {
	_, err := uuid.Parse(eventId)
	if err != nil {
		return nil, errors.Wrap(err, "not proper uuid")
	}
	// var formatedUUID pgtype.UUID
	// formatedUUID.Set(eventUUID)

	selectquery := `
		select
			id,
			name,
			description,
			created_by,
			ST_X(location_date) AS location_lat,
			ST_Y(location_date) AS location_lon,
			ST_M(location_date) AS time,
			organizer,
			upvote,
			downvote,
			created_at
		from event
		where id = $1;
	`
	// Execute the query
	row := rp.db.QueryRowxContext(ctx, selectquery, eventId)
	var eventModel model.Event
	err = row.StructScan(&eventModel)
	if err != nil {
		rp.lg.Error(selectquery)
		rp.lg.Error(eventId)
		return nil, errors.Wrap(err, "Failed to execute SQL query")
	}
	return &eventModel, nil
}
