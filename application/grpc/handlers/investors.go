package handlers

import (
	"context"
	"log/slog"

	"connectrpc.com/connect"

	depositsv1 "github.com/iainvm/deposits/application/grpc/gen/deposits/v1"
	"github.com/iainvm/deposits/internal/investors"
)

type InvestorsService interface {
	Onboard(ctx context.Context, investor *investors.Investor) error
}

type InvestorsHandler struct {
	log              *slog.Logger
	investorsService InvestorsService
}

func NewInvestorsHandler(log *slog.Logger, service InvestorsService) *InvestorsHandler {
	return &InvestorsHandler{
		log:              log,
		investorsService: service,
	}
}

func (h *InvestorsHandler) Onboard(ctx context.Context, req *connect.Request[depositsv1.OnboardRequest]) (*connect.Response[depositsv1.OnboardResponse], error) {
	h.log.With("header", req.Header()).With("request", req.Msg).Info("Onboard Called")

	// Create domain model
	investor, err := investors.NewInvestor(req.Msg.Investor.Name)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	// Onboard
	err = h.investorsService.Onboard(ctx, investor)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	// Create response
	res := connect.NewResponse(&depositsv1.OnboardResponse{
		Investor: &depositsv1.Investor{
			Id:   investor.Id.String(),
			Name: investor.Name.String(),
		},
	})
	res.Header().Set("Investor-Version", "v1")
	return res, nil
}
