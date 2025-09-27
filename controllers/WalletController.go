package controllers

import (
	"bougette-backend/common"
	"bougette-backend/dtos"
	"bougette-backend/models"
	"bougette-backend/services"
	"bougette-backend/validation"

	"github.com/labstack/echo/v4"
)

type WalletController struct {
	WalletService *services.WalletService
}

func NewWalletController(walletService *services.WalletService) *WalletController {
	return &WalletController{WalletService: walletService}
}

func (w *WalletController) CreateWallet(ctx echo.Context) error {
	userID, ok := ctx.Get("user").(uint)
	if !ok {
		return common.SendInternalServerErrorResponse(ctx, "User authentication required")
	}

	request := new(dtos.CreateWalletRequestDTO)
	if err := ctx.Bind(request); err != nil {
		return common.SendBadRequestResponse(ctx, err.Error())
	}

	if err := validation.ValidateStruct(request); err != nil {
		validationErrors := validation.FormatValidationErrors(err)
		return common.SendFailedValidationResponse(ctx, validationErrors)
	}

	wallet := models.Wallet{
		UserID:  userID,
		Name:    request.Name,
		Balance: request.Balance,
	}

	exists, err := w.WalletService.CheckWalletExitWithUserIDAndName(userID, request.Name)
	if err != nil {
		return common.SendInternalServerErrorResponse(ctx, "Failed to check if wallet exists")
	}
	if exists {
		return common.SendBadRequestResponse(ctx, "Wallet with this user id and name already exists")
	}

	if err := w.WalletService.CreateWallet(&wallet); err != nil {
		return common.SendInternalServerErrorResponse(ctx, "Wallet could not be created")
	}

	return common.SendSuccessResponse(ctx, "Wallet created successfully", wallet)
}
