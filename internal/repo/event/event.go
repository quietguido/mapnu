package event

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/quietguido/mapnu/internal/repo/event/model"
	"github.com/quietguido/mapnu/pkg/assert"
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
	insertQuery := rp.builder.
		Insert(eventTable).Columns(
		"name",
		"description",
		"created_by",
		"location_date",
		"organizer",
		"upvote",
		"downvote",
	).Values(
		createEvent.Name,
		createEvent.Description,
		createEvent.CreatedBy,
		createPointM(
			createEvent.Location_lon,
			createEvent.Location_lat,
			createEvent.Time,
		),
		createEvent.Organizer,
		createEvent.Upvote,
		createEvent.Downvote,
	)

	sql, args, err := insertQuery.ToSql()
	assert.IsNil(err, "Failed to build SQL query")

	result, err := rp.db.ExecContext(ctx, sql, args...)
	assert.IsNil(err, "Failed to execute SQL query")
	// if err != nil {
	// 	return errors.Wrap(err, "Failed to execute SQL query")
	// }

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
	err := row.StructScan(&eventModel)
	if err != nil {
		rp.lg.Error(selectquery)
		rp.lg.Error(eventId)
		return nil, errors.Wrap(err, "Failed to execute SQL query")
	}
	return &eventModel, nil
}

func (rp *repository) GetMapForQuadrant(ctx context.Context, mapQuery model.GetMapQueryParams) ([]model.Event, error) {
	// SQL query with placeholders
	selectquery := `
	SELECT
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
	FROM event
	WHERE
		ST_Within(
			location_date,
			ST_SetSRID(
				ST_MakeEnvelope($1, $2, $3, $4, 4326), -- Dynamic spatial (lat/lon) bounds
				4326
			)
		)
		AND ST_M(location_date) BETWEEN $5 AND $6; -- Dynamic time filter
	`

	// Prepare arguments for query placeholders
	args := []interface{}{
		mapQuery.FirstQuadLon,    // $1
		mapQuery.FirstQuadLat,    // $2
		mapQuery.SecondQuadLon,   // $3
		mapQuery.SecondQuadLat,   // $4
		mapQuery.FromTime.Unix(), // $5
		mapQuery.ToTime.Unix(),   // $6
	}

	// Execute the query
	rows, err := rp.db.QueryxContext(ctx, selectquery, args...)
	if err != nil {
		rp.lg.Error("Failed to execute GetMapForQuadrant query", zap.Error(err))
		return nil, errors.Wrap(err, "Failed to execute SQL query")
	}
	defer rows.Close()

	// Parse the results
	var events []model.Event
	for rows.Next() {
		var event model.Event
		err := rows.StructScan(&event)
		if err != nil {
			rp.lg.Error("Failed to scan row", zap.Error(err))
			return nil, errors.Wrap(err, "Failed to scan row")
		}
		events = append(events, event)
	}

	// Check for any errors during iteration
	if err := rows.Err(); err != nil {
		rp.lg.Error("Row iteration error", zap.Error(err))
		return nil, errors.Wrap(err, "Row iteration error")
	}

	return events, nil
}
