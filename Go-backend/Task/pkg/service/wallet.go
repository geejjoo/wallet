package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/geejjoo/task"
	"github.com/geejjoo/task/pkg/repository"
	"time"
)

type WalletService struct {
	repo repository.Repository
}

func NewWalletService(repo repository.Repository) Wallet {
	return &WalletService{repo: repo}
}

func (w *WalletService) UpdateWallet(updateWallet *task.UpdateWallet) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := w.repo.UpdateWallet(ctx, *updateWallet)

	switch {
	case errors.Is(err, repository.FromIdError):
		return FromIdError
	case errors.Is(err, repository.ToIdError):
		return ToIdError
	case errors.Is(err, repository.BalanceError):
		return BalanceError
	}
	if err != nil {
		return fmt.Errorf("wallet to error %s", err)
	}

	err = w.repo.CreateHistory(*updateWallet)
	if err != nil {
		return err
	}

	return nil
}

func (w *WalletService) GetWallet(id string) (string, float64, error) {
	return w.repo.GetWallet(id)
}

func (w *WalletService) CreateWallet() (string, float64, error) {
	return w.repo.CreateWallet()
}

func (w *WalletService) GetHistory(id string) ([]task.History, error) {
	if _, _, err := w.repo.GetWallet(id); err != nil {
		return nil, WalletNotFoundError
	}
	history, err := w.repo.History.GetHistory(id)

	if err != nil {
		return history, DatabaseError
	}

	if history == nil {
		return []task.History{}, nil
	}

	return history, nil
}
