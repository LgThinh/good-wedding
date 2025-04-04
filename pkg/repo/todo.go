package repo

import (
	"context"
	"github.com/google/uuid"
	"good-template-go/pkg/errors"
	"good-template-go/pkg/model"
	"good-template-go/pkg/utils"
	"good-template-go/pkg/utils/logger"
	"gorm.io/gorm"
)

// TodoRepo is a struct that contains the database implementation for truck entity
type TodoRepo struct {
	DB *gorm.DB
}

func (r *TodoRepo) DBWithTimeout(ctx context.Context) (*gorm.DB, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(ctx, generalQueryTimeout)
	return r.DB.WithContext(ctx), cancel
}

func NewRepoTodo(todoRepo *gorm.DB) TodoRepoInterface {
	return &TodoRepo{
		DB: todoRepo,
	}
}

type TodoRepoInterface interface {
	DBWithTimeout(ctx context.Context) (*gorm.DB, context.CancelFunc)
	Create(tx *gorm.DB, ob *model.Todo) (*model.Todo, error)
	Get(tx *gorm.DB, id uuid.UUID) (*model.Todo, error)
	Filter(tx *gorm.DB, f *model.TodoFilter) (*model.TodoFilterResult, error)
	Update(tx *gorm.DB, id, updaterID uuid.UUID, ob *model.TodoUpdateRequest) (*model.Todo, error)
	Delete(tx *gorm.DB, id, updaterID uuid.UUID) error
	GetOneFlexible(tx *gorm.DB, field string, value interface{}) (*model.Todo, error)
}

func (r *TodoRepo) Create(tx *gorm.DB, ob *model.Todo) (*model.Todo, error) {
	log := logger.WithTag("TodoRepo|Create")
	err := tx.Create(&ob).Error
	if err != nil {
		logger.LogError(log, err, "Error when get create")
		appErr := errors.FeAppError("Không tạo được", errors.UnknownError)
		return nil, appErr
	}

	return ob, nil
}

func (r *TodoRepo) Get(tx *gorm.DB, id uuid.UUID) (*model.Todo, error) {
	log := logger.WithTag("TodoRepo|Get")
	var todo model.Todo
	err := tx.Where("id = ?", id).First(&todo).Error
	if err != nil {
		logger.LogError(log, err, "Error when get information")
		appErr := errors.FeAppError("Không tìm thấy thông tin", errors.UnknownError)
		return nil, appErr
	}

	return &todo, nil
}

func (r *TodoRepo) Filter(tx *gorm.DB, f *model.TodoFilter) (*model.TodoFilterResult, error) {
	log := logger.WithTag("TodoRepo|Filter")

	tx = tx.Model(&model.Todo{})

	if f.FromDate != nil {
		fromDate := utils.ConvertUnixMilliToTime(*f.FromDate)
		tx = tx.Where("created_at >= ?", fromDate)
	}

	if f.ToDate != nil {
		toDate := utils.ConvertUnixMilliToTime(*f.ToDate)
		tx = tx.Where("created_at <= ?", toDate)
	}

	if f.CreatorID != nil {
		tx = tx.Where("creator_id = ?", *f.CreatorID)
	}

	if f.Name != nil {
		tx = tx.Where("name = ?", *f.Name)
	}

	if f.Key != nil {
		tx = tx.Where("key = ?", *f.Key)
	}

	if f.IsActive != nil {
		tx = tx.Where("is_active = ?", *f.IsActive)
	}

	if f.Code != nil {
		tx = tx.Where("code = ?", *f.Code)
	}

	result := &model.TodoFilterResult{
		Filter:  f,
		Records: []*model.Todo{},
	}

	f.Pager.SortableFields = []string{"id", "created_at", "updated_at", "name"}

	pager := f.Pager

	tx = pager.DoQuery(&result.Records, tx)
	if tx.Error != nil {
		logger.LogError(log, tx.Error, "Error when get list")
		appErr := errors.FeAppError("Không tìm thấy danh sách", errors.NotFound)
		return nil, appErr
	}

	if result.Records == nil {
		result.Records = []*model.Todo{}
	}

	return result, nil
}

func (r *TodoRepo) Update(tx *gorm.DB, id, updaterID uuid.UUID, ob *model.TodoUpdateRequest) (*model.Todo, error) {
	log := logger.WithTag("TodoRepo|Update")

	var updatedTodo model.Todo

	if err := tx.Model(&updatedTodo).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"updater_id":  updaterID,
			"name":        *ob.Name,
			"key":         *ob.Key,
			"is_active":   *ob.IsActive,
			"code":        *ob.Code,
			"description": *ob.Description,
		}).
		First(&updatedTodo, id).Error; err != nil {
		logger.LogError(log, err, "Error when updating")
		return nil, errors.FeAppError("Không cập nhập được", errors.UnknownError)
	}

	return &updatedTodo, nil
}

func (r *TodoRepo) Delete(tx *gorm.DB, id, updaterID uuid.UUID) error {
	log := logger.WithTag("TodoRepo|Delete")

	if err := tx.Delete(model.Todo{}, id).Error; err != nil {
		logger.LogError(log, err, "Error when delete")
		err = errors.FeAppError("Không xóa được", errors.UnknownError)
		return err
	}

	return nil
}

func (r *TodoRepo) GetOneFlexible(tx *gorm.DB, field string, value interface{}) (*model.Todo, error) {
	log := logger.WithTag("TodoRepo|GetOneFlexible")

	ob := &model.Todo{}

	if err := tx.Where(field+" = ? ", value).First(&ob).Error; err != nil {
		logger.LogError(log, err, "Error when get")
		err = errors.FeAppError("Không tìm thấy", errors.UnknownError)
		return nil, err
	}

	return ob, nil
}
