package repository

import (
	"github.com/geejjoo/task"
	"github.com/jmoiron/sqlx"
	"strings"
	"time"
)

type HistoryDB struct {
	db            *sqlx.DB
	queryReplacer *strings.Replacer
}

func NewHistoryDB(db *sqlx.DB, tableName string) History {
	return &HistoryDB{
		db:            db,
		queryReplacer: strings.NewReplacer("{table}", tableName),
	}
}

func (w *HistoryDB) CreateHistory(wallet task.UpdateWallet) error {
	query := w.queryReplacer.Replace("INSERT INTO {table} (from_id, to_id, amount, time) values ($1, $2, $3, $4)")
	row := w.db.QueryRow(query, wallet.FromID, wallet.ToID, wallet.Amount, time.Now())
	if row.Err() != nil {
		return DatabaseError
	}
	return nil
}

func (w *HistoryDB) GetHistory(id string) ([]task.History, error) {
	var history []task.History
	query := w.queryReplacer.Replace("SELECT from_id, to_id, amount, time FROM {table} WHERE from_id=$1 or to_id=$2")
	err := w.db.Select(&history, query, id, id)
	if err != nil {
		return nil, err
	}
	return history, nil
}
