package handlers

import (
	"context"
	"log/slog"

	"connectrpc.com/connect"

	depositsv1 "github.com/iainvm/deposits/application/grpc/gen/deposits/v1"
	"github.com/iainvm/deposits/internal/investors"
)

type InvestorsHandler struct {
	investorsService *investors.Service
	log              *slog.Logger
}

func NewInvestorsHandler(service *investors.Service, log *slog.Logger) *InvestorsHandler {
	return &InvestorsHandler{
		investorsService: service,
		log:              log,
	}
}

func (h *InvestorsHandler) Onboard(ctx context.Context, req *connect.Request[depositsv1.OnboardRequest]) (*connect.Response[depositsv1.OnboardResponse], error) {
	h.log.With("header", req.Header()).With("request", req.Msg).Info("Onboard Called")

	//
	investor, err := h.investorsService.Onboard(ctx, req.Msg.Investor.Name)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
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
