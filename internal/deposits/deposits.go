package deposits

import (
	"errors"

	"github.com/google/uuid"
)

var (
	ErrIdGeneration            = errors.New("failed to generate id")
	ErrWrapperTypeExistsInPot  = errors.New("pot already contains wrapper type")
	ErrNominalAmountNegative   = errors.New("nominal amount cannot be negative value")
	ErrAllocatedAmountNegative = errors.New("allocated amount cannot be negative value")
)

type DepositId string

func newDepositId() (DepositId, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return "", errors.Join(ErrIdGeneration, err)
	}

	return DepositId(id.String()), nil
}

func (id DepositId) String() string {
	return string(id)
}

type Deposit struct {
	Id   DepositId
	Pots []*Pot
}

func New() (*Deposit, error) {
	// Generate Id
	id, err := newDepositId()
	if err != nil {
		return nil, err
	}

	// Create Deposit
	return &Deposit{
		Id: id,
	}, nil
}

func (deposit *Deposit) AddPot(pot *Pot) {
	deposit.Pots = append(deposit.Pots, pot)
}
