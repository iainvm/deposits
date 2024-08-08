package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"

	"github.com/sethvargo/go-envconfig"

	"github.com/iainvm/deposits/common/postgres"
	"github.com/iainvm/deposits/internal/deposits"
	depositsStore "github.com/iainvm/deposits/internal/deposits/postgres"
	"github.com/iainvm/deposits/internal/investors"
	investorsStore "github.com/iainvm/deposits/internal/investors/postgres"
)

type DBConfig struct {
	Host     string `env:"HOST, default=localhost"`
	Port     string `env:"PORT, default=5432"`
	User     string `env:"USER, default=postgres"`
	Password string `env:"PASSWORD, default=postgres"`
	Name     string `env:"NAME, default=postgres"`
}

type Config struct {
	Port     string   `env:"PORT, default=8080"`
	DBConfig DBConfig `env:", prefix=DB_"`
}

func main() {
	// Logger
	ctx := context.Background()
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	slog.SetDefault(logger)

	// Parse Env Vars
	var config Config
	err := envconfig.Process(ctx, &config)
	if err != nil {
		logger.With("error", err).Error("failed to parse server configuration")
		panic(fmt.Errorf("failed to parse server configuration: %w", err))
	}
	logger.Debug("Config Processed", "Config", config)

	// DB
	dataSource := postgres.NewDataSource(
		config.DBConfig.Host,
		config.DBConfig.Port,
		config.DBConfig.User,
		config.DBConfig.Password,
		config.DBConfig.Name,
		false,
	)
	db, err := postgres.Connect(dataSource)
	if err != nil {
		logger.With("error", err).Error("failed to connect to DB")
		panic(fmt.Errorf("failed to connect to DB: %w", err))
	}
	logger.With("host", config.DBConfig.Host).With("port", config.DBConfig.Port).Info("Connected to DB")

	investorsService := investors.NewService(
		investorsStore.NewStore(db),
	)

	depositsService := deposits.NewService(
		depositsStore.NewStore(db),
	)

	PlayThrough(investorsService, depositsService)
}

func PlayThrough(investorsService *investors.Service, depositsService *deposits.Service) {
	ctx := context.Background()

	// Create investor data
	investor, err := investors.NewInvestor("Iain")
	if err != nil {
		panic(err)
	}

	// Onboard
	err = investorsService.Onboard(ctx, investor)
	if err != nil {
		panic(err)
	}

	// Create Deposit
	deposit, err := deposits.NewDeposit()
	if err != nil {
		panic(err)
	}

	// Add Pot A to deposit
	potA, err := deposits.NewPot("Pot A")
	if err != nil {
		panic(err)
	}
	deposit.AddPot(potA)

	// Add GIA account to Pot A
	accountGIA, err := deposits.NewAccount(
		deposits.WrapperTypeGIA,
		10_000,
	)
	if err != nil {
		panic(err)
	}
	err = potA.AddAccount(accountGIA)
	if err != nil {
		panic(err)
	}

	// Add ISA account to Pot A
	accountISA, err := deposits.NewAccount(
		deposits.WrapperTypeISA,
		20_000,
	)
	if err != nil {
		panic(err)
	}
	err = potA.AddAccount(accountISA)
	if err != nil {
		panic(err)
	}

	// Add SIPP account to Pot A
	accountSIPP, err := deposits.NewAccount(
		deposits.WrapperTypeSIPP,
		50_000,
	)
	if err != nil {
		panic(err)
	}
	err = potA.AddAccount(accountSIPP)
	if err != nil {
		panic(err)
	}

	// Add Pot B to deposit
	potB, err := deposits.NewPot("Pot B")
	if err != nil {
		panic(err)
	}
	deposit.AddPot(potB)

	// Add GIA account to Pot B
	accountGIA2, err := deposits.NewAccount(
		deposits.WrapperTypeGIA,
		20_000,
	)
	if err != nil {
		panic(err)
	}
	err = potB.AddAccount(accountGIA2)
	if err != nil {
		panic(err)
	}

	// Now the data should be setup like shown in the PDF
	data, err := json.MarshalIndent(deposit, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(data))

	// We can create and store the deposit
	err = depositsService.Create(ctx, investor.Id, deposit)
	if err != nil {
		panic(err)
	}

	// We can create receipts
	receipt, err := deposits.NewReceipt(5_000)
	if err != nil {
		panic(err)
	}

	// And receive them against an account
	err = depositsService.ReceiveReceipt(ctx, accountGIA.Id, receipt)
	if err != nil {
		panic(err)
	}

	// Get updated deposit data
	deposit, err = depositsService.Get(ctx, deposit.Id)
	if err != nil {
		panic(err)
	}

	// View the deposit
	data, err = json.MarshalIndent(deposit, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(data))

	// GIA Accounts can go over
	receipt, err = deposits.NewReceipt(100_000)
	if err != nil {
		panic(err)
	}
	err = depositsService.ReceiveReceipt(ctx, accountGIA.Id, receipt)
	if err != nil {
		panic(err)
	}

	// ISA Accounts can't go over
	receipt, err = deposits.NewReceipt(100_000)
	if err != nil {
		panic(err)
	}
	err = depositsService.ReceiveReceipt(ctx, accountISA.Id, receipt)
	if err == nil {
		panic("ISA Account was allowed to go over the limit")
	}

	// SIPP Accounts can't go over
	receipt, err = deposits.NewReceipt(100_000)
	if err != nil {
		panic(err)
	}
	err = depositsService.ReceiveReceipt(ctx, accountSIPP.Id, receipt)
	if err == nil {
		panic("SIPP Account was allowed to go over the limit")
	}

	// Final view of the deposit
	data, err = json.MarshalIndent(deposit, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(data))
}
