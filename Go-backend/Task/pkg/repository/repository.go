package repository

import (
	"github.com/jmoiron/sqlx"
)

type Repository struct {
	Wallet
	History
}

func NewRepository(
	db *sqlx.DB,
	walletTableName string,
	historyTableName string,
) *Repository {
	return &Repository{
		NewWalletDB(db, walletTableName),
		NewHistoryDB(db, historyTableName),
	}
}
