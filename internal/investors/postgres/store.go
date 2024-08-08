package store

import (
	"context"
	"errors"

	"github.com/iainvm/deposits/internal/investors"
	"github.com/jmoiron/sqlx"
)

var ErrCreationFailed = errors.New("failed to create investor")

type Store struct {
	db *sqlx.DB
}

func NewStore(db *sqlx.DB) Store {
	return Store{
		db: db,
	}
}

type InvestorRow struct {
	Id   string `db:"id"`
	Name string `db:"name"`
}

// SaveInvestor saves the given investor to the connected database
func (store Store) SaveInvestor(ctx context.Context, investor *investors.Investor) error {
	// Define query separately for easy editting
	const query = `--sql
	INSERT INTO investors (id, name)
	VALUES (:id, :name)
	`
	// Create Row
	row := InvestorRow{
		Id:   investor.Id.String(),
		Name: investor.Name.String(),
	}

	// Execute query
	_, err := store.db.NamedExecContext(
		ctx,
		query,
		row,
	)
	if err != nil {
		return errors.Join(ErrCreationFailed, err)
	}

	return nil
}
