package investors

import "context"

type Repository interface {
	SaveInvestor(ctx context.Context, investor *Investor) error
}

type Service struct {
	repository Repository
}

func NewService(store Repository) *Service {
	return &Service{
		repository: store,
	}
}

func (service Service) Onboard(ctx context.Context, investor *Investor) error {
	// Store data
	err := service.repository.SaveInvestor(ctx, investor)
	if err != nil {
		return err
	}

	return nil
}
