package controllers

import (
	"bougette-backend/common"
	"bougette-backend/services"

	"github.com/labstack/echo/v4"
)

// UploadRequest represents the request body for upload
type UploadRequest struct {
	Filename string `json:"filename" validate:"required"`
}

type UploadController struct {
	uploadService *services.UploadService
}

func NewUploadController(uploadService *services.UploadService) *UploadController {
	return &UploadController{
		uploadService: uploadService,
	}
}

func (uc *UploadController) GeneratePresignedUploadURL(c echo.Context) error {
	var req UploadRequest
	if err := c.Bind(&req); err != nil {
		return common.SendBadRequestResponse(c, "Invalid request body")
	}

	if req.Filename == "" {
		return common.SendBadRequestResponse(c, "Filename is required")
	}

	uploadURL, key, err := uc.uploadService.GeneratePresignedUploadURL(req.Filename)
	if err != nil {
		return common.SendInternalServerErrorResponse(c, "Failed to generate presigned upload URL: "+err.Error())
	}

	return common.SendSuccessResponse(c, "Presigned upload URL generated successfully", map[string]interface{}{
		"upload_url": uploadURL,
		"key":        key,
		"expires_in": 3600,
	})
}

func (uc *UploadController) GeneratePresignedDownloadURL(c echo.Context) error {
	key := c.Param("key")
	if key == "" {
		return common.SendBadRequestResponse(c, "Key parameter is required")
	}

	url, err := uc.uploadService.GeneratePresignedDownloadURL(key)
	if err != nil {
		return common.SendInternalServerErrorResponse(c, "Failed to generate presigned URL: "+err.Error())
	}

	return common.SendSuccessResponse(c, "Presigned URL generated successfully", map[string]interface{}{
		"url":        url,
		"expires_in": 3600,
	})
}

func (uc *UploadController) DeleteFile(c echo.Context) error {
	key := c.Param("key")
	if key == "" {
		return common.SendBadRequestResponse(c, "Key parameter is required")
	}

	err := uc.uploadService.DeleteFile(key)
	if err != nil {
		return common.SendInternalServerErrorResponse(c, "Failed to delete file: "+err.Error())
	}

	return common.SendSuccessResponse(c, "File deleted successfully", map[string]string{
		"key": key,
	})
}
