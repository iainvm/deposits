package deposits

import (
	"context"

	"github.com/iainvm/deposits/internal/investors"
)

type Repository interface {
	SaveDeposit(ctx context.Context, investor_id investors.InvestorId, deposit Deposit) error
	SavePot(ctx context.Context, deposit_id DepositId, pot Pot) error
	SaveAccount(ctx context.Context, pot_id PotId, account Account) error
}

type Service struct {
	repository Repository
}

func NewService(repository Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (service *Service) Create(ctx context.Context, investor_id investors.InvestorId, deposit *Deposit) error {
	// Save Deposit
	err := service.repository.SaveDeposit(ctx, investor_id, *deposit)
	if err != nil {
		return err
	}

	// Save Pots
	for _, pot := range deposit.Pots {
		err := service.repository.SavePot(ctx, deposit.Id, *pot)
		if err != nil {
			return err
		}

		// Save Accounts
		for _, account := range pot.Accounts {
			err := service.repository.SaveAccount(ctx, pot.Id, *account)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
