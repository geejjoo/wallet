package repository

import (
	"context"
	"database/sql"
	"github.com/geejjoo/task"
	"github.com/jmoiron/sqlx"
	"strings"
)

type WalletDB struct {
	db            *sqlx.DB
	queryReplacer *strings.Replacer
}

func NewWalletDB(db *sqlx.DB, tableName string) Wallet {
	return &WalletDB{
		db:            db,
		queryReplacer: strings.NewReplacer("{table}", tableName),
	}
}

func (w *WalletDB) CreateWallet() (string, float64, error) {
	var wallet task.Wallet

	query := w.queryReplacer.Replace("INSERT INTO {table} (balance) values ($1) RETURNING id, balance")
	row := w.db.QueryRow(query, initialBalance)
	if err := row.Scan(&wallet.ID, &wallet.Balance); err != nil {
		return "", 0, DatabaseError
	}
	return wallet.ID, wallet.Balance, nil
}

func (w *WalletDB) GetWallet(id string) (string, float64, error) {
	var wallet task.Wallet

	query := w.queryReplacer.Replace("SELECT * FROM {table} WHERE ID=$1")
	row := w.db.QueryRow(query, id)
	if err := row.Scan(&wallet.ID, &wallet.Balance); err != nil {
		return "", 0, FromIdError
	}
	return wallet.ID, wallet.Balance, nil
}

func (w *WalletDB) UpdateWallet(ctx context.Context, sendWallet task.UpdateWallet) error {
	tx, err := w.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return err
	}

	var wallet task.Wallet

	query := w.queryReplacer.Replace("SELECT * FROM {table} WHERE ID=$1")
	row := tx.QueryRow(query, sendWallet.FromID)
	if err := row.Scan(&wallet.ID, &wallet.Balance); err != nil {
		return FromIdError
	}

	query = w.queryReplacer.Replace("SELECT * FROM {table} WHERE ID=$1")
	row = tx.QueryRow(query, sendWallet.ToID)
	if err := row.Scan(&wallet.ID, &wallet.Balance); err != nil {
		return ToIdError
	}

	if wallet.Balance < sendWallet.Amount {
		return BalanceError
	}

	query = w.queryReplacer.Replace(`
UPDATE {table} 
SET balance = CASE 
	WHEN ID = $1 THEN balance - $3
	WHEN ID = $2 THEN balance + $3
END
WHERE ID IN ($1, $2)
	`)
	_, err = tx.Exec(query, sendWallet.FromID, sendWallet.ToID, sendWallet.Amount)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}

	return nil
}
