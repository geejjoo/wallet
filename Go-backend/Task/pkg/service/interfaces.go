package service

import (
	"github.com/geejjoo/task"
)

type Wallet interface {
	CreateWallet() (string, float64, error)
	GetWallet(id string) (string, float64, error)
	UpdateWallet(updateWallet *task.UpdateWallet) error
	GetHistory(id string) ([]task.History, error)
}
