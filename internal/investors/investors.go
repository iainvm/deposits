package investors

import (
	"errors"

	"github.com/google/uuid"
)

var (
	ErrIdGeneration    = errors.New("failed to generate id")
	ErrInvalidInvestor = errors.New("invalid investor")
	ErrInvalidId       = errors.New("invalid id")
	ErrInvalidName     = errors.New("invalid name")
	ErrBlankName       = errors.New("blank name given")
)

type Investor struct {
	Id   InvestorId
	Name Name
}

func New(name string) (*Investor, error) {
	id, err := newInvestorId()
	if err != nil {
		return nil, nil
	}

	investorsName, err := NewName(name)
	if err != nil {
		return nil, errors.Join(ErrInvalidInvestor, err)
	}

	investor := &Investor{
		Id:   id,
		Name: investorsName,
	}

	return investor, nil
}

type InvestorId string

func newInvestorId() (InvestorId, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return "", errors.Join(ErrIdGeneration, err)
	}

	return InvestorId(id.String()), nil
}

func ParseInvestorId(id string) (InvestorId, error) {
	_, err := uuid.Parse(id)
	if err != nil {
		return "", err
	}

	return InvestorId(id), nil
}

func (id InvestorId) String() string {
	return string(id)
}

type Name string

func NewName(name string) (Name, error) {
	if name == "" {
		return "", ErrBlankName
	}

	return Name(name), nil
}

func (name Name) String() string {
	return string(name)
}
