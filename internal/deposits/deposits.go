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

type Deposit struct {
	Id   DepositId
	Pots []*Pot
}

func newDepositId() (DepositId, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return "", errors.Join(ErrIdGeneration, err)
	}

	return DepositId(id.String()), nil
}

func ParseDepositId(id string) (DepositId, error) {
	_, err := uuid.Parse(id)
	if err != nil {
		return "", err
	}

	return DepositId(id), nil
}

func (id DepositId) String() string {
	return string(id)
}

// NewDeposit creates a new Deposit with a new Id
func NewDeposit() (*Deposit, error) {
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

// ParseDeposit parses the given data into a Deposit type, ensuring it's valid data
func ParseDeposit(id string) (*Deposit, error) {
	depositId, err := ParseDepositId(id)
	if err != nil {
		return nil, err
	}

	return &Deposit{
		Id: depositId,
	}, nil
}

func (deposit *Deposit) AddPot(pot *Pot) {
	deposit.Pots = append(deposit.Pots, pot)
}
