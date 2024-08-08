package store

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

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

type FullDeposit struct {
	Id                   string `db:"id"`
	InvestorId           string `db:"investor_id"`
	PotId                string `db:"pots_id"`
	PotName              string `db:"pots_name"`
	AccountId            string `db:"account_id"`
	AccountWrapperType   int    `db:"account_wrapper_type"`
	AccountNominalAmount int64  `db:"account_nominal_amount"`
}

func (store Store) GetFullDeposit(ctx context.Context, depositId deposits.DepositId) (*deposits.Deposit, error) {
	const query = `--sql
	SELECT d.id AS "id",
		d.investor_id AS "investor_id",
		p.id AS "pots_id",
		p.name AS "pots_name",
		a.id AS "account_id",
		a.wrapper_type AS "account_wrapper_type",
		a.nominal_amount AS "account_nominal_amount"
	FROM deposits d
	JOIN pots p ON d.id = p.deposit_id
	JOIN accounts a ON p.id = a.pot_id
	WHERE d.id = $1
	`

	rows := []FullDeposit{}

	err := store.db.Select(&rows, query, depositId.String())
	if err != nil {
		return nil, err
	}

	slog.Info("Row Count", "len", len(rows))
	deposit, err := createDomainDeposit(rows)
	if err != nil {
		return nil, err
	}
	return deposit, nil
}

func createDomainDeposit(rows []FullDeposit) (*deposits.Deposit, error) {
	if len(rows) == 0 {
		return nil, fmt.Errorf("no deposit data")
	}

	// Create the deposit
	depositId := rows[0].Id
	deposit, err := deposits.ParseDeposit(depositId)
	if err != nil {
		return nil, err
	}

	potIndexes := map[string]int{}
	for _, row := range rows {
		// Check if pot exists
		var pot *deposits.Pot
		potIndex, ok := potIndexes[row.PotId]

		// Create pot if doesn't exist
		if !ok {
			pot, err = deposits.ParsePot(row.PotId, row.PotName)
			if err != nil {
				return nil, err
			}

			// Index to find next row
			potIndexes[row.PotId] = len(deposit.Pots)
			deposit.AddPot(pot)
		} else {
			// Get pot if it exists
			pot = deposit.Pots[potIndex]
		}

		account, err := deposits.ParseAccount(row.AccountId, row.AccountWrapperType, row.AccountNominalAmount)
		if err != nil {
			return nil, err
		}
		err = pot.AddAccount(account)
		if err != nil {
			return nil, err
		}
	}

	return deposit, nil
}

func (store Store) GetDeposit(ctx context.Context, depositId deposits.DepositId) (*deposits.Deposit, error) {
	const query = `--sql
	SELECT *
	FROM deposits
	WHERE id=$1
	`

	row := DepositRow{}
	err := store.db.Get(&row, query, depositId.String())
	if err != nil {
		return nil, err
	}

	deposit, err := deposits.ParseDeposit(row.Id)
	if err != nil {
		return nil, err
	}

	return deposit, nil
}

func (store Store) SaveDeposit(ctx context.Context, investorId investors.InvestorId, deposit deposits.Deposit) error {
	// Define query separately for easy editting
	const query = `--sql
	INSERT INTO deposits (id, investor_id)
	VALUES (:id, :investor_id)
	`

	// Create Row
	row := DepositRow{
		Id:         deposit.Id.String(),
		InvestorId: investorId.String(),
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
	WrapperType   int    `db:"wrapper_type"`
	NominalAmount int64  `db:"nominal_amount"`
}

func (store Store) SaveAccount(ctx context.Context, potId deposits.PotId, account deposits.Account) error {
	// Define query separately for easy editting
	const query = `--sql
	INSERT INTO accounts (id, pot_id, wrapper_type, nominal_amount)
	VALUES (:id, :pot_id, :wrapper_type, :nominal_amount)
	`

	// Create Row
	row := AccountRow{
		Id:            account.Id.String(),
		PotId:         potId.String(),
		WrapperType:   account.WrapperType.Int(),
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
