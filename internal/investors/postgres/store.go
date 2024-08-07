package store

import (
	"context"
	"errors"

	"github.com/iainvm/deposits/internal/investors"
	"github.com/jmoiron/sqlx"
)

var ErrCreationFailed = errors.New("failed to create investor")

const Table = "investors"

type Store struct {
	db *sqlx.DB
}

func NewStore(db *sqlx.DB) Store {
	return Store{
		db: db,
	}
}

type Investor struct {
	Id   string `db:"id"`
	Name string `db:"name"`
}

// Create
func (store Store) Create(ctx context.Context, investor investors.Investor) error {
	// Define query separately for easy editting
	const query = `--sql
	INSERT INTO investors (id, name)
	VALUES (:id, :name)
	`

	// Execute query
	_, err := store.db.NamedExecContext(
		ctx,
		query,
		investor,
	)
	if err != nil {
		return errors.Join(ErrCreationFailed, err)
	}

	return nil
}
