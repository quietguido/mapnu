package event

import (
	"context"
	"fmt"
	"strconv"

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

/*
CREATE TABLE event_2025_03_15
PARTITION OF event
FOR VALUES FROM ('2025-03-15 00:00:00') TO ('2025-03-16 00:00:00');
*/

func (rp *repository) CreateEvent(ctx context.Context, createEvent model.CreateEvent) (int, error) {
	partition_event_table := getPartition(createEvent.StartDate)

	insertQuery := rp.builder.
		Insert(partition_event_table).Columns(
		"name",
		"description",
		"created_by",
		"location",
		"start_date",
		"organizer",
	).Values(
		createEvent.Name,
		createEvent.Description,
		createEvent.CreatedBy,
		sq.Expr("ST_SetSRID(ST_Point(?, ?), 4326)", createEvent.Location_lon, createEvent.Location_lat),
		createEvent.StartDate,
		createEvent.Organizer,
	).Suffix("RETURNING event_id")

	sql, args, err := insertQuery.ToSql()
	assert.IsNil(err, "Failed to build SQL query")

	var eventID int
	err = rp.db.QueryRowContext(ctx, sql, args...).Scan(&eventID)
	if err != nil {
		rp.lg.Warn(sql)
		return 0, errors.Wrap(err, "Failed to execute SQL query")
	}

	return eventID, nil
}

func (rp *repository) GetEventById(ctx context.Context, eventId int) (*model.Event, error) {
	selectquery := `
		SELECT
			event_id,
			name,
			description,
			created_by,
			ST_X(location) AS location_lon, -- Ensure longitude is first (PostGIS standard)
			ST_Y(location) AS location_lat, -- Latitude second
			start_date,
			organizer,
			upvote,
			downvote,
			created_at
		FROM event
		WHERE event_id = $1;
	`

	// Execute the query
	row := rp.db.QueryRowxContext(ctx, selectquery, eventId)
	var eventModel model.Event
	err := row.StructScan(&eventModel)
	if err != nil {
		rp.lg.Error("SQL Query Failed:", zap.String("query", selectquery))
		rp.lg.Error("Event ID:", zap.String("event_id", strconv.Itoa(eventId)))
		return nil, errors.Wrap(err, "Failed to execute SQL query")
	}
	return &eventModel, nil
}

func (rp *repository) GetMapForQuadrant(ctx context.Context, mapQuery model.GetMapQueryParams) ([]model.Event, error) {
	// ✅ Dynamically get partition table name
	partitionTable := getPartition(mapQuery.Date) // Format: "event_YYYY_MM_DD"

	selectquery := fmt.Sprintf(`
		SELECT
			event_id,
			name,
			description,
			created_by,
			ST_X(location) AS location_lon, -- Longitude
			ST_Y(location) AS location_lat, -- Latitude
			start_date,
			organizer,
			upvote,
			downvote,
			created_at
		FROM %s
		WHERE
			ST_Within(
				location,
				ST_SetSRID(
					ST_MakeEnvelope($1, $2, $3, $4, 4326), -- Bounding box for spatial query
					4326
				)
			)
	`, partitionTable) // ✅ Inject partition table name

	args := []interface{}{
		mapQuery.FirstQuadLon,  // $1 - Min Longitude
		mapQuery.FirstQuadLat,  // $2 - Min Latitude
		mapQuery.SecondQuadLon, // $3 - Max Longitude
		mapQuery.SecondQuadLat, // $4 - Max Latitude
	}

	rows, err := rp.db.QueryxContext(ctx, selectquery, args...)
	if err != nil {
		rp.lg.Error("Failed to execute GetMapForQuadrant query", zap.Error(err))
		return nil, errors.Wrap(err, "Failed to execute SQL query")
	}
	defer rows.Close()

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

	if err := rows.Err(); err != nil {
		rp.lg.Error("Row iteration error", zap.Error(err))
		return nil, errors.Wrap(err, "Row iteration error")
	}

	return events, nil
}
