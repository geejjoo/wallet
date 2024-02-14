package repository

import (
	"context"
	"errors"
	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/geejjoo/task"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateWallet(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmockDB error %s", err)
	}
	defer mockDB.Close()
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
	repo := NewWalletDB(sqlxDB, "Wallet")

	t.Run("Test create wallet in DB", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "balance"}).
			AddRow("TEST", 100.00)
		mock.ExpectQuery("INSERT").WillReturnRows(rows)

		id, balance, err := repo.CreateWallet()
		assert.NoError(t, err)
		assert.Equal(t, 100.0, balance)
		assert.Equal(t, "TEST", id)

	})

	t.Run("Test bad create wallet in DB", func(t *testing.T) {
		mock.ExpectQuery("SELECT").WillReturnError(err)

		_, _, err := repo.CreateWallet()
		assert.ErrorIs(t, err, DatabaseError)

	})

}

func TestGetWallet(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmockDB error %s", err)
	}
	defer mockDB.Close()
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
	repo := NewWalletDB(sqlxDB, "Wallet")

	t.Run("Test get wallet from DB", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "balance"}).
			AddRow("TEST", 100.00)
		mock.ExpectQuery("SELECT").WillReturnRows(rows)

		id, balance, err := repo.GetWallet("TEST")
		assert.NoError(t, err)
		assert.Equal(t, 100.00, balance)
		assert.Equal(t, "TEST", id)

	})

	t.Run("Test bad get wallet from DB", func(t *testing.T) {
		mock.ExpectQuery("SELECT").WillReturnError(err)

		_, _, err := repo.GetWallet("TEST")
		assert.ErrorIs(t, err, FromIdError)

	})

}

func TestUpdateWallet(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmockDB error %s", err)
	}
	defer mockDB.Close()
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
	repo := NewWalletDB(sqlxDB, "Wallet")

	t.Run("Test successful update", func(t *testing.T) {
		// Ожидание начала транзакции
		mock.ExpectBegin()

		// Ожидание выполнения запросов SELECT
		mock.ExpectQuery("SELECT").
			WithArgs("sender_id").
			WillReturnRows(sqlmock.NewRows([]string{"ID", "Balance"}).AddRow("sender_id", 500.0))
		mock.ExpectQuery("SELECT").
			WithArgs("receiver_id").
			WillReturnRows(sqlmock.NewRows([]string{"ID", "Balance"}).AddRow("receiver_id", 200.0))

		// Ожидание выполнения запроса UPDATE
		mock.ExpectExec(`
        UPDATE Wallet`).
			WithArgs("sender_id", "receiver_id", 100.0).
			WillReturnResult(sqlmock.NewResult(1, 1))

		// Ожидание коммита транзакции
		mock.ExpectCommit()

		err := repo.UpdateWallet(context.TODO(), task.UpdateWallet{
			FromID: "sender_id",
			ToID:   "receiver_id",
			Amount: 100.0,
		})
		assert.NoError(t, err)
	})

	t.Run("Test failed update", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(`
			UPDATE WALLET`).WithArgs("sender_id", "receiver_id", 100.0).
			WillReturnError(errors.New("update error"))
		mock.ExpectRollback()

		err := repo.UpdateWallet(context.TODO(), task.UpdateWallet{
			FromID: "sender_id",
			ToID:   "receiver_id",
			Amount: 100.0,
		})
		assert.Error(t, err)
	})
}
