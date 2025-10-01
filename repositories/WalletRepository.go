package repositories

import (
	"bougette-backend/models"

	"gorm.io/gorm"
)

type WalletRepository struct {
	db *gorm.DB
}

func NewWalletRepository(db *gorm.DB) *WalletRepository {
	return &WalletRepository{db: db}
}

func (w *WalletRepository) WalletsList(userID uint) ([]models.Wallet, error) {
	var wallets []models.Wallet
	if err := w.db.Where("user_id = ?", userID).Find(&wallets).Error; err != nil {
		return nil, err
	}
	return wallets, nil
}

func (w *WalletRepository) CreateWallet(wallet *models.Wallet) error {
	return w.db.Create(wallet).Error
}

func (w *WalletRepository) CheckWalletExitWithUserIDAndName(userID uint, name string) (*models.Wallet, error) {
	var wallet models.Wallet
	if err := w.db.Where("user_id = ? AND name = ?", userID, name).First(&wallet).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &wallet, nil
}
