package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"good-template-go/pkg/errors"
	"good-template-go/pkg/model"
	"good-template-go/pkg/model/paging"
	"good-template-go/pkg/service"
	"net/http"
)

// TodoHandler is a struct that contains the Todo router
type TodoHandler struct {
	todoService service.TodoServiceInterface
}

func NewTodoHandler(todoService service.TodoServiceInterface) *TodoHandler {
	return &TodoHandler{todoService: todoService}
}

// Create godoc
// @Summary	Create new TODO
// @Tags		TODO
// @Security   Authorization
// @Security   User ID
// @Param	todo	body	model.CreateTodoRequest	true	"New TODO"
// @Router		/todo/create [post]
func (h *TodoHandler) Create(ctx *gin.Context) {
	var (
		request model.CreateTodoRequest
	)

	adminID, exists := ctx.Get("id")
	if !exists {
		appErr := errors.FeAppError("Không tìm thấy ID trong token", errors.ValidationError)
		_ = ctx.Error(appErr)
		return
	}

	adminUUID, ok := adminID.(uuid.UUID)
	if !ok {
		appErr := errors.FeAppError("ID không đúng định dạng UUID", errors.ValidationError)
		_ = ctx.Error(appErr)
		return
	}

	if err := ctx.ShouldBind(&request); err != nil {
		appErr := errors.FeAppError(errors.VnValidationErrorMessage, errors.ValidationError)
		_ = ctx.Error(appErr)
		return
	}

	rs, err := h.todoService.Create(ctx, request, adminUUID)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, rs)
}

// Get godoc
// @Summary	Get TODO
// @Tags		TODO
// @Security   Authorization
// @Security   User ID
// @Param		id	path		string	true	"id todo"
// @Router		/todo/get-one/{id} [get]
func (h *TodoHandler) Get(ctx *gin.Context) {
	IdString := ctx.Param("id")
	id, err := uuid.Parse(IdString)
	if err != nil {
		appErr := errors.FeAppError("ID không hợp lệ", errors.BadRequest)
		_ = ctx.Error(appErr)
		return
	}

	rs, err := h.todoService.Get(ctx, id)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, rs)
}

// List godoc
// @Summary	List TODO
// @Tags		TODO
// @Security   Authorization
// @Security   User ID
// @Param		page_size	query		int		true	"size per page"
// @Param		page		query		int		true	"page number (> 0)"
// @Param		sort		query		string	false	"sort"
// @Router		/todo/get-list/ [get]
func (h *TodoHandler) List(ctx *gin.Context) {
	var (
		request model.TodoFilterRequest
	)
	err := ctx.ShouldBindQuery(&request)
	if err != nil {
		err = errors.FeAppError(errors.VnValidationErrorMessage, errors.BadRequest)
		_ = ctx.Error(err)
		return
	}

	filter := &model.TodoFilter{
		TodoFilterRequest: request,
		Pager:             paging.NewPagerWithGinCtx(ctx),
	}

	rs, err := h.todoService.Filter(ctx, filter)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, paging.NewBodyPaginated(rs.Records, rs.Filter.Pager))
}

// Update godoc
// @Summary	Update TODO
// @Tags		TODO
// @Security   Authorization
// @Security   User ID
// @Param		id	path		string				true	"id"
// @Param		todo	body		model.TodoRequest	true	"Update todo"
// @Router		/todo/update/{id} [put]
func (h *TodoHandler) Update(ctx *gin.Context) {
	var (
		request model.TodoUpdateRequest
	)
	if err := ctx.ShouldBind(&request); err != nil {
		err = errors.FeAppError(errors.VnValidationErrorMessage, errors.BadRequest)
		_ = ctx.Error(err)
		return
	}

	updaterID, exists := ctx.Get("id")
	if !exists {
		appErr := errors.FeAppError("Không tìm thấy ID trong token", errors.ValidationError)
		_ = ctx.Error(appErr)
		return
	}

	updaterUUID, ok := updaterID.(uuid.UUID)
	if !ok {
		appErr := errors.FeAppError("ID không đúng định dạng UUID", errors.ValidationError)
		_ = ctx.Error(appErr)
		return
	}

	IdString := ctx.Param("id")
	id, err := uuid.Parse(IdString)
	if err != nil {
		appErr := errors.FeAppError("ID không hợp lệ", errors.BadRequest)
		_ = ctx.Error(appErr)
		return
	}

	// Update
	rs, err := h.todoService.Update(ctx, id, updaterUUID, request)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, rs)
}

// Delete godoc
// @Summary	Delete TODO
// @Tags		TODO
// @Security   Authorization
// @Security   User ID
// @Param		id	path		string				true	"id"
// @Router		/todo/delete/{id} [delete]
func (h *TodoHandler) Delete(ctx *gin.Context) {
	updaterID, exists := ctx.Get("id")
	if !exists {
		appErr := errors.FeAppError("Không tìm thấy ID trong token", errors.ValidationError)
		_ = ctx.Error(appErr)
		return
	}

	updaterUUID, ok := updaterID.(uuid.UUID)
	if !ok {
		appErr := errors.FeAppError("ID không đúng định dạng UUID", errors.ValidationError)
		_ = ctx.Error(appErr)
		return
	}

	IdString := ctx.Param("id")
	id, err := uuid.Parse(IdString)
	if err != nil {
		appErr := errors.FeAppError("ID không hợp lệ", errors.BadRequest)
		_ = ctx.Error(appErr)
		return
	}

	// Delete
	err = h.todoService.Delete(ctx, id, updaterUUID)
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, "Delete Success")
}
