package repository

import (
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/geejjoo/task"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestCreateHistory(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmockDB error %s", err)
	}
	defer mockDB.Close()
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
	repo := NewHistoryDB(sqlxDB, "History")

	t.Run("Test successful history creation", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"from_id", "to_id", "amount", time.DateTime}).
			AddRow("from_id", "to_id", 100.0, time.Now())
		mock.ExpectQuery("INSERT INTO ").WillReturnRows(rows)

		err := repo.CreateHistory(task.UpdateWallet{
			FromID: "from_id",
			ToID:   "to_id",
			Amount: 100.0,
		})
		assert.NoError(t, err)
	})

	t.Run("Test failed history creation", func(t *testing.T) {
		mock.ExpectQuery("INSERT INTO HISTORY").
			WithArgs("from_id", "to_id", 100.0, time.Now()).
			WillReturnError(errors.New("insert error"))

		err := repo.CreateHistory(task.UpdateWallet{
			FromID: "from_id",
			ToID:   "to_id",
			Amount: 100.0,
		})
		assert.Error(t, err)
	})
}

func TestGetHistory(t *testing.T) {
	// Setup
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmockDB error %s", err)
	}
	defer mockDB.Close()
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
	repo := NewHistoryDB(sqlxDB, "History")

	t.Run("Test successful retrieval of history", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"from_id", "to_id", "amount", "time"}).
			AddRow("from_id_1", "to_id_1", 100.0, time.Now()).
			AddRow("from_id_2", "to_id_2", 200.0, time.Now())
		mock.ExpectQuery("SELECT").
			WillReturnRows(rows)

		history, err := repo.GetHistory("id")
		assert.NoError(t, err)
		assert.Len(t, history, 2)
	})

	t.Run("Test failed retrieval of history", func(t *testing.T) {
		mock.ExpectQuery("SELECT").
			WillReturnError(errors.New("select error"))

		history, err := repo.GetHistory("id")
		assert.Error(t, err)
		assert.Nil(t, history)
	})
}

func TestNewHistoryDB(t *testing.T) {
	db := &sqlx.DB{}
	historyDB := NewHistoryDB(db, "History")

	assert.NotNil(t, historyDB)

}
