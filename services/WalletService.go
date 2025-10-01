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

func (s *WalletService) WalletsList(userID uint) ([]models.Wallet, error) {
	return s.WalletRepos.WalletsList(userID)
}

func (s *WalletService) CheckWalletExitWithUserIDAndName(userID uint, name string) (bool, error) {
	wallet, err := s.WalletRepos.CheckWalletExitWithUserIDAndName(userID, name)
	if err != nil {
		return false, err
	}
	return wallet != nil, nil
}

func (s *WalletService) GenerateDefaultWallet(userID uint) ([]*models.Wallet, error) {
	wallets := []string{"Groceries", "Transportation", "Entertainment", "Health", "Education", "Other"}
	var walletsCreated []*models.Wallet

	for _, walletName := range wallets {
		walletExits, err := s.WalletRepos.CheckWalletExitWithUserIDAndName(userID, walletName)
		if err != nil {
			return nil, err
		}

		if walletExits != nil {
			walletsCreated = append(walletsCreated, walletExits)
			continue
		}

		newWallet := &models.Wallet{
			UserID:  userID,
			Name:    walletName,
			Balance: 0,
		}

		if err := s.CreateWallet(newWallet); err != nil {
			return nil, err
		}

		walletsCreated = append(walletsCreated, newWallet)
	}

	return walletsCreated, nil
}
