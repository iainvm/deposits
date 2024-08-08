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

func (wrapperType WrapperType) Int() int {
	return int(wrapperType)
}

type AccountId string

func newAccountId() (AccountId, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return "", errors.Join(ErrIdGeneration, err)
	}

	return AccountId(id.String()), nil
}

func ParseAccountId(id string) (AccountId, error) {
	_, err := uuid.Parse(id)
	if err != nil {
		return "", err
	}

	return AccountId(id), nil
}

func (id AccountId) String() string {
	return string(id)
}

type NominalAmount int

type AllocatedAmount int

func NewAccount(wrapperType WrapperType, nominalAmount int64) (*Account, error) {
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

func ParseAccount(id string, wrapperType int, nominalAmount int64) (*Account, error) {
	accountId, err := ParseAccountId(id)
	if err != nil {
		return nil, err
	}

	accountNominalAmount, err := NewNominalAmount(nominalAmount)
	if err != nil {
		return nil, err
	}

	accountWrapperType := WrapperType(wrapperType)

	account := &Account{
		Id:            accountId,
		WrapperType:   accountWrapperType,
		NominalAmount: accountNominalAmount,
	}

	return account, nil
}

func NewNominalAmount(amount int64) (NominalAmount, error) {
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
