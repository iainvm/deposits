package investors

import "context"

type Repository interface {
	Create(ctx context.Context, investor Investor) error
}

type Service struct {
	repository Repository
}

func NewService(store Repository) *Service {
	return &Service{
		repository: store,
	}
}

func (service Service) Create(ctx context.Context, name string) (Investor, error) {

	// Create Data, ensures valid
	inv, err := New(name)
	if err != nil {
		return Investor{}, err
	}

	// Store data
	err = service.repository.Create(ctx, inv)
	if err != nil {
		return Investor{}, err
	}

	return inv, nil
}
