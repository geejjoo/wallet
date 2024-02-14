package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
)

const (
	initialBalance = 100.00
)

func NewDB(cfg *Config) (*sqlx.DB, error) {
	db, err := sqlx.Open(cfg.DriverName, fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.DBName, cfg.Password, cfg.SSLMode))
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}
