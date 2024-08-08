package deposits

import (
	"errors"

	"github.com/google/uuid"
)

type Receipt struct {
	Id              ReceiptId
	AccountId       AccountId
	AllocatedAmount AllocatedAmount
}

type AllocatedAmount int64

type ReceiptId string

func NewReceipt(allocatedAmount int64) (*Receipt, error) {
	id, err := newReceiptId()
	if err != nil {
		return nil, err
	}

	amount, err := NewAllocatedAmount(allocatedAmount)
	if err != nil {
		return nil, err
	}
	return &Receipt{
		Id:              id,
		AllocatedAmount: amount,
	}, nil
}

func ParseReceipt(id string, allocatedAmount int64) (*Receipt, error) {
	receiptId, err := ParseReceiptId(id)
	if err != nil {
		return nil, err
	}

	receiptAllocatedAmount, err := NewAllocatedAmount(allocatedAmount)
	if err != nil {
		return nil, err
	}

	receipt := &Receipt{
		Id:              receiptId,
		AllocatedAmount: receiptAllocatedAmount,
	}

	return receipt, nil
}

func NewAllocatedAmount(amount int64) (AllocatedAmount, error) {
	if amount < 0 {
		return 0, ErrAllocatedAmountNegative
	}

	return AllocatedAmount(amount), nil
}

func (allocatedAmount AllocatedAmount) Int64() int64 {
	return int64(allocatedAmount)
}

func newReceiptId() (ReceiptId, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return "", errors.Join(ErrIdGeneration, err)
	}

	return ReceiptId(id.String()), nil
}

func ParseReceiptId(id string) (ReceiptId, error) {
	_, err := uuid.Parse(id)
	if err != nil {
		return "", err
	}

	return ReceiptId(id), nil
}

func (id ReceiptId) String() string {
	return string(id)
}
