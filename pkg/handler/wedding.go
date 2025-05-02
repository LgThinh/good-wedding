package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"good-wedding/pkg/errors"
	"good-wedding/pkg/model"
	"good-wedding/pkg/model/paging"
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

func (h *WeddingHandler) Comment(ctx *gin.Context) {
	var (
		request model.CommentRequest
	)
	if err := ctx.ShouldBind(&request); err != nil {
		appErr := errors.FeAppError(errors.VnValidationErrorMessage, errors.ValidationError)
		_ = ctx.Error(appErr)
		return
	}

	rs, err := h.weddingService.Comment(ctx, request)
	if err != nil {
		_ = ctx.Error(err)
	}

	ctx.JSON(http.StatusOK, rs)
}

func (h *WeddingHandler) WeddingWish(ctx *gin.Context) {
	var (
		request model.WeddingWishRequest
	)
	if err := ctx.ShouldBind(&request); err != nil {
		appErr := errors.FeAppError(errors.VnValidationErrorMessage, errors.ValidationError)
		_ = ctx.Error(appErr)
		return
	}

	rs, err := h.weddingService.WeddingWish(ctx, request)
	if err != nil {
		_ = ctx.Error(err)
	}

	ctx.JSON(http.StatusOK, rs)
}

func (h *WeddingHandler) ListComment(ctx *gin.Context) {
	var (
		request model.CommentFilterRequest
	)
	err := ctx.ShouldBindQuery(&request)
	if err != nil {
		err = errors.FeAppError(errors.VnValidationErrorMessage, errors.BadRequest)
		_ = ctx.Error(err)
		return
	}

	filter := &model.CommentFilter{
		CommentFilterRequest: request,
		Pager:                paging.NewPagerWithGinCtx(ctx),
	}

	rs, err := h.weddingService.CommentFilter(ctx, filter)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, paging.NewBodyPaginated(rs.Records, rs.Filter.Pager))
}

func (h *WeddingHandler) ListWeddingWish(ctx *gin.Context) {
	var (
		request model.WeddingWishFilterRequest
	)
	err := ctx.ShouldBindQuery(&request)
	if err != nil {
		err = errors.FeAppError(errors.VnValidationErrorMessage, errors.BadRequest)
		_ = ctx.Error(err)
		return
	}

	filter := &model.WeddingWishFilter{
		WeddingWishFilterRequest: request,
		Pager:                    paging.NewPagerWithGinCtx(ctx),
	}

	rs, err := h.weddingService.WeddingWishFilter(ctx, filter)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, paging.NewBodyPaginated(rs.Records, rs.Filter.Pager))
}

func (h *WeddingHandler) ListUser(ctx *gin.Context) {
	var (
		request model.UserFilterRequest
	)
	err := ctx.ShouldBindQuery(&request)
	if err != nil {
		err = errors.FeAppError(errors.VnValidationErrorMessage, errors.BadRequest)
		_ = ctx.Error(err)
		return
	}

	filter := &model.UserFilter{
		UserFilterRequest: request,
		Pager:             paging.NewPagerWithGinCtx(ctx),
	}

	rs, err := h.weddingService.UserFilter(ctx, filter)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, paging.NewBodyPaginated(rs.Records, rs.Filter.Pager))
}

func (h *WeddingHandler) ListMedia(ctx *gin.Context) {
	var (
		request model.ObjectMediaFilterRequest
	)
	err := ctx.ShouldBindQuery(&request)
	if err != nil {
		err = errors.FeAppError(errors.VnValidationErrorMessage, errors.BadRequest)
		_ = ctx.Error(err)
		return
	}

	filter := &model.ObjectMediaFilter{
		ObjectMediaFilterRequest: request,
		Pager:                    paging.NewPagerWithGinCtx(ctx),
	}

	rs, err := h.weddingService.ObjectMediaFilter(ctx, filter)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, paging.NewBodyPaginated(rs.Records, rs.Filter.Pager))
}
