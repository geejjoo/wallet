package service

import (
	"context"
	"github.com/geejjoo/task"
	"github.com/geejjoo/task/pkg/repository"
	"github.com/stretchr/testify/assert"
	"testing"
)

type MockWallet struct {
}

type MockHistory struct {
}

func (w *MockWallet) CreateWallet() (string, float64, error) {
	return "", 0, nil
}

func (w *MockWallet) GetWallet(id string) (string, float64, error) {
	return "", 0, nil
}

func (w *MockWallet) UpdateWallet(ctx context.Context, sendWallet task.UpdateWallet) error {
	switch {
	case sendWallet.FromID == "1":
		return repository.FromIdError
	case sendWallet.ToID == "1":
		return repository.ToIdError
	case sendWallet.Amount >= 111:
		return repository.BalanceError
	}
	return nil
}

func (w *MockHistory) CreateHistory(sendWallet task.UpdateWallet) error {
	return nil
}

func (w *MockHistory) GetHistory(id string) ([]task.History, error) {
	return []task.History{}, nil
}

func TestUpdateWallet(t *testing.T) {

	walletRepo := repository.Repository{&MockWallet{}, &MockHistory{}}
	walletService := NewWalletService(walletRepo)

	t.Run("FromID error", func(t *testing.T) {

		err := walletService.UpdateWallet(&task.UpdateWallet{FromID: "1"})
		assert.Equal(t, FromIdError, err)
	})

	t.Run("ToID error", func(t *testing.T) {

		err := walletService.UpdateWallet(&task.UpdateWallet{ToID: "1"})
		assert.Equal(t, ToIdError, err)
	})

	t.Run("Balance error", func(t *testing.T) {

		err := walletService.UpdateWallet(&task.UpdateWallet{Amount: 200})
		assert.Equal(t, BalanceError, err)
	})

	t.Run("Update and create history success", func(t *testing.T) {

		err := walletService.UpdateWallet(&task.UpdateWallet{FromID: "12", ToID: "123", Amount: 50})
		assert.NoError(t, err)
	})
}

func TestGetWallet(t *testing.T) {
	walletRepo := repository.Repository{&MockWallet{}, &MockHistory{}}
	walletService := NewWalletService(walletRepo)

	t.Run("Get wallet check", func(t *testing.T) {

		_, _, err := walletService.GetWallet("1")
		assert.NoError(t, err)
	})

}

func TestCreateWallet(t *testing.T) {
	walletRepo := repository.Repository{&MockWallet{}, &MockHistory{}}
	walletService := NewWalletService(walletRepo)

	t.Run("Create wallet check", func(t *testing.T) {

		_, _, err := walletService.CreateWallet()
		assert.NoError(t, err)
	})

}

func TestGetHistory(t *testing.T) {
	walletRepo := repository.Repository{&MockWallet{}, &MockHistory{}}
	walletService := NewWalletService(walletRepo)

	t.Run("Get history check", func(t *testing.T) {

		_, err := walletService.GetHistory("1")
		assert.NoError(t, err)
	})

}
