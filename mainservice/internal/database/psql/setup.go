package psql

import (
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

func New(cfg Config) (*sqlx.DB, error) {
	dataSource := fmt.Sprintf(
		"user=%s password=%s host=%s port=%s dbname=%s sslmode=disable",
		cfg.User, cfg.Password, cfg.Addr, cfg.Port, cfg.DB,
	)
	conn, err := sqlx.Connect("pgx", dataSource)
	if err != nil {
		return nil, errors.Wrap(err, dataSource)
	}

	err = conn.Ping()
	if err != nil {
		return nil, errors.Wrap(err, "ping failed")
	}
	return conn, nil
}
