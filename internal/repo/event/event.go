package event

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/quietguido/mapnu/internal/repo/event/model"
)

const (
	event = "event"
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
	insertquery := rp.builder.
		Insert(event).Columns(
		"name",
		"description",
		"created_by",
		"location_date",
		"organizer",
		"upvote",
		"downvote",
		"created_at",
	).Values(
		createEvent.Name,
		createEvent.Description,
		createEvent.CreatedBy,
		createST_GeomFromText(createEvent.Location_lon, createEvent.Location_lat, createEvent.Time),
		createEvent.Organizer,
		createEvent.Upvote,
		createEvent.Downvote,
		createEvent.CreatedAt,
	)

	sql, args, err := insertquery.ToSql()
	if err != nil {
		return errors.Wrap(err, "Failed to create sql query")
	}

	_, err = rp.db.ExecContext(ctx, sql, args)

	if err != nil {
		return errors.Wrap(err, "Failed to run sql query")
	}

	return nil
}

// func (rp *repository) CreateEvent() error {
// 	selectquery := rp.builder.
// 		Select(
// 			"id",
// 			"name",
// 			"description",
// 			"created_by",
// 			"location_lat",
// 			"location_lon",
// 			"time",
// 			"organizer",
// 			"upvote",
// 			"downvote",
// 			"created_at",
// 		).
// 		From(event)

// 	return nil
// }
