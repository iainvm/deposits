package deposits

import (
	"context"

	"github.com/iainvm/deposits/internal/investors"
)

type Repository interface {
	SaveDeposit(ctx context.Context, investorId investors.InvestorId, deposit Deposit) error
	SavePot(ctx context.Context, depositId DepositId, pot Pot) error
	SaveAccount(ctx context.Context, potId PotId, account Account) error
	SaveReceipt(ctx context.Context, accountId AccountId, receipt Receipt) error
	GetFullDeposit(ctx context.Context, depositId DepositId) (*Deposit, error)
	GetAccount(ctx context.Context, accountId AccountId) (*Account, error)
	UpdateAccount(ctx context.Context, account Account) error
}

type Service struct {
	repository Repository
}

func NewService(repository Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (service *Service) ReceiveReceipt(ctx context.Context, accountId AccountId, receipt *Receipt) error {

	// Get Account
	account, err := service.repository.GetAccount(ctx, accountId)
	if err != nil {
		return err
	}

	// Validate we can add the receipt to the account
	err = account.AddReceipt(receipt)
	if err != nil {
		return err
	}

	// Save the receipt
	err = service.repository.SaveReceipt(ctx, account.Id, *receipt)
	if err != nil {
		return err
	}

	// Update the account
	err = service.repository.UpdateAccount(ctx, *account)
	if err != nil {
		return err
	}

	return nil
}

func (service *Service) Get(ctx context.Context, id DepositId) (*Deposit, error) {
	deposit, err := service.repository.GetFullDeposit(ctx, id)
	if err != nil {
		return nil, err
	}
	return deposit, nil
}

func (service *Service) Create(ctx context.Context, investorId investors.InvestorId, deposit *Deposit) error {
	// Save Deposit
	err := service.repository.SaveDeposit(ctx, investorId, *deposit)
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
