package repository

import "errors"

var DatabaseError = errors.New("Database error")

var FromIdError = errors.New("Incorrect wallet ID")

var ToIdError = errors.New("Incorrect target ID")

var BalanceError = errors.New("Not enough money")
