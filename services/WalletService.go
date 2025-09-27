package services

import (
	"bougette-backend/models"
	"bougette-backend/repositories"
)

type WalletService struct {
	WalletRepos *repositories.WalletRepository
}

func NewWalletService(walletRepo *repositories.WalletRepository) *WalletService {
	return &WalletService{WalletRepos: walletRepo}
}

func (s *WalletService) CreateWallet(wallet *models.Wallet) error {
	return s.WalletRepos.CreateWallet(wallet)
}

func (s *WalletService) CheckWalletExitWithUserIDAndName(userID uint, name string) (bool, error) {
	wallet, err := s.WalletRepos.CheckWalletExitWithUserIDAndName(userID, name)
	if err != nil {
		return false, err
	}
	return wallet != nil, nil
}
