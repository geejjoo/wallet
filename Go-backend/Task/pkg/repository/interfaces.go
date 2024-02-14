package repository

import (
	"context"
	"github.com/geejjoo/task"
)

type Wallet interface {
	CreateWallet() (string, float64, error)
	GetWallet(id string) (string, float64, error)
	UpdateWallet(ctx context.Context, wallet task.UpdateWallet) error
}

type History interface {
	CreateHistory(wallet task.UpdateWallet) error
	GetHistory(id string) ([]task.History, error)
}
