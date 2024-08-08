package deposits

import (
	"errors"

	"github.com/google/uuid"
)

type WrapperType int

const (
	WrapperTypeGIA WrapperType = iota + 1 //Bump number so that you can't set UNSPECIFIED account type
	WrapperTypeISA
	WrapperTypeSIPP
)

type Account struct {
	Id              AccountId
	WrapperType     WrapperType
	NominalAmount   NominalAmount
	AllocatedAmount AllocatedAmount
}

type AccountId string

func newAccountId() (AccountId, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return "", errors.Join(ErrIdGeneration, err)
	}

	return AccountId(id.String()), nil
}

func (id AccountId) String() string {
	return string(id)
}

type NominalAmount int

type AllocatedAmount int

func NewAccount(wrapperType WrapperType, nominalAmount int) (*Account, error) {
	// Generate Id
	id, err := newAccountId()
	if err != nil {
		return nil, err
	}

	accountNominalAmount, err := NewNominalAmount(nominalAmount)
	if err != nil {
		return nil, err
	}

	// Create Account
	return &Account{
		Id:            id,
		WrapperType:   wrapperType,
		NominalAmount: accountNominalAmount,
	}, nil
}

func NewNominalAmount(amount int) (NominalAmount, error) {
	if amount < 0 {
		return 0, ErrNominalAmountNegative
	}

	return NominalAmount(amount), nil
}

func (nominalAmount NominalAmount) Int64() int64 {
	return int64(nominalAmount)
}

func NewAllocatedAmount(amount int) (AllocatedAmount, error) {
	if amount < 0 {
		return 0, ErrAllocatedAmountNegative
	}

	return AllocatedAmount(amount), nil
}

func (allocatedAmount AllocatedAmount) Int64() int64 {
	return int64(allocatedAmount)
}
