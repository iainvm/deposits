package handlers

import (
	"context"
	"log/slog"

	"connectrpc.com/connect"
	depositsv1 "github.com/iainvm/deposits/application/grpc/gen/deposits/v1"
	"github.com/iainvm/deposits/internal/deposits"
	"github.com/iainvm/deposits/internal/investors"
)

type DepositsHandler struct {
	log              *slog.Logger
	depostitsService *deposits.Service
}

func NewDepositsHandler(log *slog.Logger, depositsService *deposits.Service) *DepositsHandler {
	return &DepositsHandler{
		log:              log,
		depostitsService: depositsService,
	}
}

func (h *DepositsHandler) Get(ctx context.Context, req *connect.Request[depositsv1.GetRequest]) (*connect.Response[depositsv1.GetResponse], error) {
	h.log.With("header", req.Header()).With("request", req.Msg).Info("Get Called")

	depositId, err := deposits.ParseDepositId(req.Msg.Id)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	deposit, err := h.depostitsService.Get(ctx, depositId)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	// Create response
	res := connect.NewResponse(&depositsv1.GetResponse{
		Deposit: createResponseDeposit(*deposit),
	})
	res.Header().Set("Deposit-Version", "v1")
	return res, nil
}

func (h *DepositsHandler) Create(ctx context.Context, req *connect.Request[depositsv1.CreateRequest]) (*connect.Response[depositsv1.CreateResponse], error) {
	h.log.With("header", req.Header()).With("request", req.Msg).Info("Create Called")

	// Create Domain Model
	deposit, err := createDomainDeposit(req.Msg.Deposit)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	investorId, err := investors.ParseInvestorId(req.Msg.InvestorId)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	// Onboard
	err = h.depostitsService.Create(ctx, investorId, deposit)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	// Create response
	res := connect.NewResponse(&depositsv1.CreateResponse{
		Deposit: createResponseDeposit(*deposit),
	})
	res.Header().Set("Deposit-Version", "v1")
	return res, nil
}

func createDomainDeposit(reqDeposit *depositsv1.Deposit) (*deposits.Deposit, error) {
	deposit, err := deposits.New()
	if err != nil {
		return nil, err
	}

	// Add Pots
	for _, reqPot := range reqDeposit.Pots {
		pot, err := deposits.NewPot(reqPot.Name)
		if err != nil {
			return nil, err
		}

		// Add Accounts
		for _, reqAccount := range reqPot.Accounts {
			// Get Wrapper Type
			wrapperType := deposits.WrapperType(reqAccount.WrapperType)
			account, err := deposits.NewAccount(wrapperType, reqAccount.NominalAmount)
			if err != nil {
				return nil, err
			}

			pot.AddAccount(account)
		}

		deposit.AddPot(pot)
	}

	return deposit, nil
}

func createResponseDeposit(deposit deposits.Deposit) *depositsv1.Deposit {

	// Create deposit
	response := &depositsv1.Deposit{
		Id:   deposit.Id.String(),
		Pots: []*depositsv1.Pot{},
	}

	// Attach pots
	for _, pot := range deposit.Pots {
		responsePot := &depositsv1.Pot{
			Id:       pot.Id.String(),
			Name:     pot.Name.String(),
			Accounts: []*depositsv1.Account{},
		}

		// Attach Accounts
		for _, account := range pot.Accounts {
			responseAccount := &depositsv1.Account{
				Id:              account.Id.String(),
				WrapperType:     depositsv1.WrapperType(account.WrapperType),
				NominalAmount:   account.NominalAmount.Int64(),
				AllocatedAmount: account.AllocatedAmount.Int64(),
			}

			responsePot.Accounts = append(responsePot.Accounts, responseAccount)
		}

		response.Pots = append(response.Pots, responsePot)
	}

	return response
}
