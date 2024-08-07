package investor

import (
	"errors"
	"log/slog"

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
	Id   Id
	Name Name
}

func New(name string) (Investor, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return Investor{}, errors.Join(ErrIdGeneration, err)
	}

	_name, err := NewName(name)
	if err != nil {
		return Investor{}, errors.Join(ErrInvalidInvestor, err)
	}

	investor := Investor{
		Id:   Id(id.String()),
		Name: _name,
	}

	// TODO: Call to store

	return investor, nil
}

type Id string

func NewId(id string) (Id, error) {
	_, err := uuid.Parse(id)
	if err != nil {
		slog.With("id", id).Error("invalid id")
		return "", errors.Join(ErrInvalidId, err)
	}

	return Id(id), nil
}

func (id Id) String() string {
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
