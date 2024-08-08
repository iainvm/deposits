package deposits

import (
	"errors"

	"github.com/google/uuid"
)

type Pot struct {
	Id       PotId
	Name     PotName
	Accounts []*Account
}

type PotId string

func newPotId() (PotId, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return "", errors.Join(ErrIdGeneration, err)
	}

	return PotId(id.String()), nil
}

func (id PotId) String() string {
	return string(id)
}

type PotName string

func NewPotName(name string) (PotName, error) {
	return PotName(name), nil
}

func (potName PotName) String() string {
	return string(potName)
}

func NewPot(name string) (*Pot, error) {
	// Generate Id
	id, err := newPotId()
	if err != nil {
		return nil, err
	}

	potName, err := NewPotName(name)
	if err != nil {
		return nil, err
	}

	// Create Pot
	return &Pot{
		Id:   id,
		Name: potName,
	}, nil
}

func (pot *Pot) AddAccount(account *Account) error {
	// Pots can contain only 1 of each wrapper type
	for _, potAccount := range pot.Accounts {
		if potAccount.WrapperType == account.WrapperType {
			return ErrWrapperTypeExistsInPot
		}
	}

	// Append account to pot
	pot.Accounts = append(pot.Accounts, account)
	return nil
}
