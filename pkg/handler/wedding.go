package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"good-wedding/pkg/errors"
	"good-wedding/pkg/service"
	"net/http"
)

type WeddingHandler struct {
	weddingService service.WeddingServiceInterface
}

func NewWeddingHandler(weddingService service.WeddingServiceInterface) *WeddingHandler {
	return &WeddingHandler{weddingService: weddingService}
}

func (h *WeddingHandler) UploadImage(ctx *gin.Context) {
	file, err := ctx.FormFile("image")
	if err != nil {
		appErr := errors.FeAppError(errors.VnValidationErrorMessage, errors.ValidationError)
		_ = ctx.Error(appErr)
		return
	}

	adminUUID, err := uuid.Parse("123e4567-e89b-12d3-a456-426614174000")
	if err != nil {
		appErr := errors.FeAppError("UUID test không hợp lệ", errors.ValidationError)
		_ = ctx.Error(appErr)
		return
	}

	rs, err := h.weddingService.UploadImageToS3(ctx, file, adminUUID)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, rs)
}

func (h *WeddingHandler) UploadVideo(ctx *gin.Context) {
	file, err := ctx.FormFile("video")
	if err != nil {
		appErr := errors.FeAppError(errors.VnValidationErrorMessage, errors.ValidationError)
		_ = ctx.Error(appErr)
		return
	}

	adminUUID, err := uuid.Parse("123e4567-e89b-12d3-a456-426614174000")
	if err != nil {
		appErr := errors.FeAppError("UUID test không hợp lệ", errors.ValidationError)
		_ = ctx.Error(appErr)
		return
	}

	rs, err := h.weddingService.UploadVideoToS3(ctx, file, adminUUID)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, rs)
}
