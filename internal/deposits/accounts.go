package deposits

import (
	"errors"

	"github.com/google/uuid"
)

type WrapperType int

var (
	ErrNominalExceeded = errors.New("nomial value exceeded")
)

const (
	WrapperTypeGIA WrapperType = iota + 1 //Bump number so that you can't set UNSPECIFIED account type
	WrapperTypeISA
	WrapperTypeSIPP
)

type Account struct {
	Id                   AccountId
	WrapperType          WrapperType
	TotalAllocatedAmount TotalAllocatedAmount
	NominalAmount        NominalAmount
	Receipts             []*Receipt
}
type AccountId string

type NominalAmount int64

type TotalAllocatedAmount int64

func (amount TotalAllocatedAmount) Int64() int64 {
	return int64(amount)
}

func NewTotalAllocatedAmount(amount int64) TotalAllocatedAmount {
	return TotalAllocatedAmount(amount)
}

func (wrapperType WrapperType) Int() int {
	return int(wrapperType)
}

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
		Id:                   id,
		WrapperType:          wrapperType,
		NominalAmount:        accountNominalAmount,
		TotalAllocatedAmount: 0,
	}, nil
}

func ParseAccount(id string, wrapperType int, nominalAmount int64, totalAllocatedAmount int64) (*Account, error) {
	accountId, err := ParseAccountId(id)
	if err != nil {
		return nil, err
	}

	accountNominalAmount, err := NewNominalAmount(nominalAmount)
	if err != nil {
		return nil, err
	}

	accountWrapperType := WrapperType(wrapperType)

	accountTotalAllocatedAmount := NewTotalAllocatedAmount(totalAllocatedAmount)

	account := &Account{
		Id:                   accountId,
		WrapperType:          accountWrapperType,
		NominalAmount:        accountNominalAmount,
		TotalAllocatedAmount: accountTotalAllocatedAmount,
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

func (account *Account) AddReceipt(receipt *Receipt) error {
	// ISA and SIPP accounts can't exceed Nominal Amount
	if account.WrapperType == WrapperTypeISA || account.WrapperType == WrapperTypeSIPP {
		if account.NominalAmount.Int64() < account.TotalAllocatedAmount.Int64()+receipt.AllocatedAmount.Int64() {
			return ErrNominalExceeded
		}
	}

	account.TotalAllocatedAmount += TotalAllocatedAmount(receipt.AllocatedAmount)
	account.Receipts = append(account.Receipts, receipt)

	return nil
}
