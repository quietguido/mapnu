package user

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/quietguido/mapnu/internal/repo/user/model"
	"github.com/quietguido/mapnu/pkg/assert"
)

const (
	userTable = "users"
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

func (rp *repository) CreateUser(ctx context.Context, newUser model.CreateUser) error {
	insertQuery := rp.builder.
		Insert(userTable).
		Columns("username", "email").
		Values(newUser.Username, newUser.Email)

	sql, args, err := insertQuery.ToSql()
	assert.IsNil(err, "Failed to build SQL query")

	result, err := rp.db.ExecContext(ctx, sql, args...)
	if err != nil {
		return errors.Wrap(err, "Failed to create user")
	}

	// Check if any rows were affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "Failed to get affected rows")
	}
	if rowsAffected == 0 {
		return errors.New("No rows were inserted")
	}

	return nil
}

// ✅ GetUserById fetches a user by ID
func (rp *repository) GetUserById(ctx context.Context, userId string) (*model.User, error) {
	getQuery := rp.builder.
		Select(
			"id",
			"username",
			"email",
			"created_at",
		).
		From("users").
		Where(sq.Eq{"id": userId})

	sql, args, err := getQuery.ToSql()
	assert.IsNil(err, "Failed to build SQL query")

	// Execute the query
	var userModel model.User
	err = rp.db.GetContext(ctx, &userModel, sql, args...)
	if err != nil {
		rp.lg.Error("Failed to execute SQL query", zap.Error(err))
		return nil, errors.Wrap(err, "Failed to fetch user")
	}

	return &userModel, nil
}
