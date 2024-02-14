package service

import (
	"github.com/geejjoo/task/pkg/repository"
)

type Service struct {
	Wallet
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		NewWalletService(*repos),
	}
}
