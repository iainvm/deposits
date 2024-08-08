package deposits

import (
	"errors"

	"github.com/google/uuid"
)

type WrapperType int

var (
	ErrNominalExceeded    = errors.New("nomial value exceeded")
	ErrNegativeAmount     = errors.New("negative amount given")
	ErrInvalidWrapperType = errors.New("invalid wrapper type given")
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

func NewTotalAllocatedAmount(amount int64) (TotalAllocatedAmount, error) {
	if amount < 0 {
		return 0, ErrNegativeAmount
	}

	return TotalAllocatedAmount(amount), nil
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

// NewAccount creates a new Account with a new Id
func NewAccount(wrapperType WrapperType, nominalAmount int64) (*Account, error) {
	// Generate Id
	id, err := newAccountId()
	if err != nil {
		return nil, err
	}

	// Wrapper Type
	err = validateWrapperType(wrapperType)
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

func validateWrapperType(wrapperType WrapperType) error {
	switch wrapperType {
	case WrapperTypeGIA, WrapperTypeISA, WrapperTypeSIPP:
		return nil
	}

	return ErrInvalidWrapperType
}

// ParseAccount parses the given data into a Account type, ensuring it's valid data
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
	err = validateWrapperType(accountWrapperType)
	if err != nil {
		return nil, err
	}

	accountTotalAllocatedAmount, err := NewTotalAllocatedAmount(totalAllocatedAmount)
	if err != nil {
		return nil, err
	}

	account := &Account{
		Id:            accountId,
		WrapperType:   accountWrapperType,
		NominalAmount: accountNominalAmount,
	}

	err = account.SetTotalAllocationAmount(accountTotalAllocatedAmount)
	if err != nil {
		return nil, err
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

// AddReceipt validates that it can allocate the receipt to the Account, then updates account information
func (account *Account) AddReceipt(receipt *Receipt) error {
	err := account.IncreaseTotalAllocationAmount(TotalAllocatedAmount(receipt.AllocatedAmount))
	if err != nil {
		return err
	}

	account.Receipts = append(account.Receipts, receipt)

	return nil
}

// IncreaseTotalAllocationAmount increases the TotalAllocatedAmount by the given amount
func (account *Account) IncreaseTotalAllocationAmount(amount TotalAllocatedAmount) error {
	newAmount := account.TotalAllocatedAmount.Int64() + amount.Int64()

	value, err := NewTotalAllocatedAmount(newAmount)
	if err != nil {
		return err
	}

	err = account.SetTotalAllocationAmount(value)
	if err != nil {
		return err
	}

	return nil
}

// SetTotalAllocationAmount sets the TotalAllocatedAmount to the given amount
func (account *Account) SetTotalAllocationAmount(amount TotalAllocatedAmount) error {

	// ISA and SIPP accounts can't exceed Nominal Amount
	if account.WrapperType == WrapperTypeISA || account.WrapperType == WrapperTypeSIPP {
		if account.NominalAmount.Int64() < amount.Int64() {
			return ErrNominalExceeded
		}
	}

	account.TotalAllocatedAmount = amount
	return nil
}
