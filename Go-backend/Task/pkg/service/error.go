package service

import "errors"

var FromIdError = errors.New("Invalid parent wallet ID")

var ToIdError = errors.New("Invalid target wallet ID")

var BalanceError = errors.New("Not enough money on wallet")

var WalletNotFoundError = errors.New("Wallet not found")

var DatabaseError = errors.New("Database error")
