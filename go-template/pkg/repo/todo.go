package repo

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"good-template-go/pkg/model"
	"good-template-go/pkg/utils"
	"gorm.io/gorm"
	"log"
)

// Repo is a struct that contains the database implementation for truck entity
type Repo struct {
	DB *gorm.DB
}

func (r *Repo) DBWithTimeout(ctx context.Context) (*gorm.DB, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(ctx, generalQueryTimeout)
	return r.DB.WithContext(ctx), cancel
}

func NewRepoTodo(repoTodo *gorm.DB) RepoTodoInterface {
	return &Repo{
		DB: repoTodo,
	}
}

type RepoTodoInterface interface {
	Create(param interface{}, ob *model.Todo) error
	Get(param interface{}, id uuid.UUID) (*model.Todo, error)
	Filter(param interface{}, f *model.TodoFilter) (*model.TodoFilterResult, error)
	Update(param interface{}, id uuid.UUID, ob *model.Todo) error
	Delete(param interface{}, id uuid.UUID) error
	GetOneFlexible(param interface{}, field string, value interface{}) (*model.Todo, error)
}

func (r *Repo) Create(param interface{}, ob *model.Todo) error {
	var (
		tx     *gorm.DB
		cancel context.CancelFunc
	)
	switch v := param.(type) {
	case context.Context:
		tx, cancel = r.DBWithTimeout(v)
		defer cancel()
	case *gorm.DB:
		tx = v
	default:
		return fmt.Errorf("invalid parameter")
	}

	return tx.Create(ob).Error
}

func (r *Repo) Get(param interface{}, id uuid.UUID) (*model.Todo, error) {
	var (
		tx     *gorm.DB
		cancel context.CancelFunc
	)
	switch v := param.(type) {
	case context.Context:
		tx, cancel = r.DBWithTimeout(v)
		defer cancel()
	case *gorm.DB:
		tx = v
	default:
		return nil, fmt.Errorf("invalid parameter")
	}

	o := &model.Todo{}

	err := tx.First(&o, id).Error
	return o, err
}

func (r *Repo) Filter(param interface{}, f *model.TodoFilter) (*model.TodoFilterResult, error) {
	var (
		tx     *gorm.DB
		cancel context.CancelFunc
	)
	switch v := param.(type) {
	case context.Context:
		tx, cancel = r.DBWithTimeout(v)
		defer cancel()
	case *gorm.DB:
		tx = v
	default:
		return nil, fmt.Errorf("invalid parameter")
	}

	tx = tx.Model(&model.Todo{})

	op := tx.Where
	tx = utils.FilterIfNotNil(f.CreatedFrom, tx, op, "created_at >= ?")
	tx = utils.FilterIfNotNil(f.CreatedTo, tx, op, "created_at <= ?")
	tx = utils.FilterIfNotNil(f.CreatorID, tx, op, "creator_id = ?")
	tx = utils.FilterIfNotNil(f.Name, tx, op, "name = ?")
	tx = utils.FilterIfNotNil(f.Key, tx, op, "key = ?")
	tx = utils.FilterIfNotNil(f.IsActive, tx, op, "is_active = ?")
	tx = utils.FilterIfNotNil(f.Code, tx, op, "code = ?")

	result := &model.TodoFilterResult{
		Filter:  f,
		Records: []*model.Todo{},
	}

	f.Pager.SortableFields = []string{"id", "created_at", "updated_at"}
	pager := result.Filter.Pager

	tx = pager.DoQuery(&result.Records, tx)
	if tx.Error != nil {
		log.Printf("error while filter todo")
	}

	return result, tx.Error
}

func (r *Repo) Update(param interface{}, id uuid.UUID, ob *model.Todo) error {
	var (
		tx     *gorm.DB
		cancel context.CancelFunc
	)
	switch v := param.(type) {
	case context.Context:
		tx, cancel = r.DBWithTimeout(v)
		defer cancel()
	case *gorm.DB:
		tx = v
	default:
		return fmt.Errorf("invalid parameter")
	}

	return tx.Where("id = ?", id).Updates(&ob).Error
}

func (r *Repo) Delete(param interface{}, id uuid.UUID) error {
	var (
		tx     *gorm.DB
		cancel context.CancelFunc
	)
	switch v := param.(type) {
	case context.Context:
		tx, cancel = r.DBWithTimeout(v)
		defer cancel()
	case *gorm.DB:
		tx = v
	default:
		return fmt.Errorf("invalid parameter")
	}

	return tx.Delete(&model.Todo{}, id).Error
}

func (r *Repo) GetOneFlexible(param interface{}, field string, value interface{}) (*model.Todo, error) {
	var (
		tx     *gorm.DB
		cancel context.CancelFunc
	)
	switch v := param.(type) {
	case context.Context:
		tx, cancel = r.DBWithTimeout(v)
		defer cancel()
	case *gorm.DB:
		tx = v
	default:
		return nil, fmt.Errorf("invalid parameter")
	}

	o := &model.Todo{}

	err := tx.Where(field+" = ? ", value).First(&o).Error
	return o, err
}
