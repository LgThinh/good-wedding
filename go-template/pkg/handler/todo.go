package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/praslar/lib/common"
	"good-template-go/pkg/model"
	"good-template-go/pkg/model/paging"
	"good-template-go/pkg/repo"
	"net/http"
)

// TodoHandler is a struct that contains the Todo service
type TodoHandler struct {
	obRepo repo.RepoTodoInterface
}

func NewTodoHandler(obRepo repo.RepoTodoInterface) *TodoHandler {
	return &TodoHandler{obRepo: obRepo}
}

// Create godoc
// @Summary	Create new TODO
// @Tags		TODO
// @Security   Authorization
// @Security   User ID
// @Param	todo	body	model.TodoRequest	true	"New TODO"
// @Router		/todo/create [post]
func (h *TodoHandler) Create(ctx *gin.Context) {
	var (
		request model.TodoRequest
	)
	userId := ctx.GetHeader("x-user-id")
	creatorId, err := uuid.Parse(userId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, "error-parse-userid")
		return
	}
	if err = ctx.ShouldBind(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, "error-parse-request")
		return
	}
	ob := &model.Todo{
		BaseModel: model.BaseModel{
			CreatorID: &creatorId,
		},
	}
	common.Sync(request, ob)

	// Create
	err = h.obRepo.Create(ctx, ob)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, "Create Failure")
		return
	}
	ctx.JSON(http.StatusOK, "Create Success")
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
		ctx.JSON(http.StatusBadRequest, "error-parse-id")
		return
	}
	// Get
	rs, err := h.obRepo.Get(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, "Fail to get data")
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
	var req model.TodoListRequest

	err := ctx.BindQuery(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "fail to get pagination"})
		return
	}

	filter := &model.TodoFilter{
		TodoListRequest: req,
		Pager:           paging.NewPagerWithGinCtx(ctx),
	}
	// List
	rs, err := h.obRepo.Filter(ctx, filter)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, "Create Failure")
		return
	}
	ctx.JSON(http.StatusOK, rs)
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
		request model.TodoRequest
	)
	userId := ctx.GetHeader("x-user-id")
	updaterId, err := uuid.Parse(userId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, "error-parse-user-id")
		return
	}
	IdString := ctx.Param("id")
	id, err := uuid.Parse(IdString)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, "error-parse-id")
		return
	}

	if err = ctx.ShouldBind(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, "error-parse-request")
		return
	}
	ob := &model.Todo{
		BaseModel: model.BaseModel{
			UpdaterID: &updaterId,
		},
	}
	common.Sync(request, ob)

	// Update
	err = h.obRepo.Update(ctx, id, ob)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, "Create Failure")
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"Updated": request,
	})
}

// Delete godoc
// @Summary	Delete TODO
// @Tags		TODO
// @Security   Authorization
// @Security   User ID
// @Param		id	path		string				true	"id"
// @Router		/todo/delete/{id} [delete]
func (h *TodoHandler) Delete(ctx *gin.Context) {
	var (
		ob *model.Todo
	)
	userId := ctx.GetHeader("x-user-id")
	updaterId, err := uuid.Parse(userId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, "error-parse-user-id")
		return
	}
	IdString := ctx.Param("id")
	id, err := uuid.Parse(IdString)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, "error-parse-id")
		return
	}
	ob.UpdaterID = &updaterId

	// Delete
	err = h.obRepo.Delete(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, "Create Failure")
		return
	}
	ctx.JSON(http.StatusOK, "Delete Success")
}

// HardDelete godoc
// @Summary	HardDelete TODO
// @Tags		TODO
// @Security   Authorization
// @Security   User ID
// @Param		id	path		string				true	"id"
// @Router		/todo/hard-delete/{id} [delete]
func (h *TodoHandler) HardDelete(ctx *gin.Context) {
	// not implement
	ctx.JSON(http.StatusOK, "Hard Delete Success")
}
