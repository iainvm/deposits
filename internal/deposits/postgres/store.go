package store

import (
	"context"
	"errors"

	"github.com/iainvm/deposits/internal/deposits"
	"github.com/iainvm/deposits/internal/investors"
	"github.com/jmoiron/sqlx"
)

var ErrSaveFailed = errors.New("failed to save deposit")

type Store struct {
	db *sqlx.DB
}

func NewStore(db *sqlx.DB) Store {
	return Store{
		db: db,
	}
}

type DepositRow struct {
	Id         string `db:"id"`
	InvestorId string `db:"investor_id"`
}

func (store Store) SaveDeposit(ctx context.Context, investor_id investors.InvestorId, deposit deposits.Deposit) error {
	// Define query separately for easy editting
	const query = `--sql
	INSERT INTO deposits (id, investor_id)
	VALUES (:id, :investor_id)
	`

	// Create Row
	row := DepositRow{
		Id:         deposit.Id.String(),
		InvestorId: investor_id.String(),
	}

	// Execute query
	_, err := store.db.NamedExecContext(
		ctx,
		query,
		row,
	)
	if err != nil {
		return errors.Join(ErrSaveFailed, err)
	}

	return nil
}

type PotRow struct {
	Id        string `db:"id"`
	DepositId string `db:"deposit_id"`
	Name      string `db:"name"`
}

func (store Store) SavePot(ctx context.Context, depositId deposits.DepositId, pot deposits.Pot) error {
	// Define query separately for easy editting
	const query = `--sql
	INSERT INTO pots (id, deposit_id, name)
	VALUES (:id, :deposit_id, :name)
	`

	// Create Row
	row := PotRow{
		Id:        pot.Id.String(),
		DepositId: depositId.String(),
		Name:      pot.Name.String(),
	}

	// Execute query
	_, err := store.db.NamedExecContext(
		ctx,
		query,
		row,
	)
	if err != nil {
		return errors.Join(ErrSaveFailed, err)
	}

	return nil
}

type AccountRow struct {
	Id            string `db:"id"`
	PotId         string `db:"pot_id"`
	NominalAmount int64  `db:"nominal_amount"`
}

func (store Store) SaveAccount(ctx context.Context, potId deposits.PotId, account deposits.Account) error {
	// Define query separately for easy editting
	const query = `--sql
	INSERT INTO accounts (id, pot_id, nominal_amount)
	VALUES (:id, :pot_id, :nominal_amount)
	`

	// Create Row
	row := AccountRow{
		Id:            account.Id.String(),
		PotId:         potId.String(),
		NominalAmount: account.NominalAmount.Int64(),
	}

	// Execute query
	_, err := store.db.NamedExecContext(
		ctx,
		query,
		row,
	)
	if err != nil {
		return errors.Join(ErrSaveFailed, err)
	}

	return nil
}
