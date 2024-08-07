package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"connectrpc.com/connect"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	investorv1 "github.com/iainvm/deposits/gen/investors/v1"         // generated by protoc-gen-go
	"github.com/iainvm/deposits/gen/investors/v1/investorsv1connect" // generated by protoc-gen-connect-go
	"github.com/iainvm/deposits/internal/investor"
)

type InvestorHandler struct {
	log *slog.Logger
}

func (h *InvestorHandler) Onboard(ctx context.Context, req *connect.Request[investorv1.OnboardRequest]) (*connect.Response[investorv1.OnboardResponse], error) {
	h.log.With("header", req.Header()).With("request", req.Msg).Info("Onboard Called")

	// Create investor
	investor, err := investor.New(req.Msg.Name)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	// Create response
	res := connect.NewResponse(&investorv1.OnboardResponse{
		Id:   investor.Id.String(),
		Name: investor.Name.String(),
	})
	res.Header().Set("Investor-Version", "v1")
	return res, nil
}

func main() {
	// Logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	slog.SetDefault(logger)

	// Init servers
	invServer := &InvestorHandler{log: logger}

	// Register handlers
	mux := http.NewServeMux()
	path, handler := investorsv1connect.NewInvestorsServiceHandler(invServer)
	mux.Handle(path, handler)

	// Listen
	http.ListenAndServe(
		"localhost:8080",
		// Use h2c so we can serve HTTP/2 without TLS.
		h2c.NewHandler(mux, &http2.Server{}),
	)
}
