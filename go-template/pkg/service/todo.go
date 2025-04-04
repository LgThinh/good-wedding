package service

import (
	"context"
	"github.com/google/uuid"
	"good-template-go/pkg/mapper"
	"good-template-go/pkg/model"
	"good-template-go/pkg/repo"
)

type TodoService struct {
	todoRepo repo.TodoRepoInterface
}

func NewTodoService(todoRepo repo.TodoRepoInterface) *TodoService {
	return &TodoService{
		todoRepo: todoRepo,
	}
}

type TodoServiceInterface interface {
	Create(ctx context.Context, req model.CreateTodoRequest, adminID uuid.UUID) (*model.CreateTodoResponse, error)
	Get(ctx context.Context, id uuid.UUID) (*model.Todo, error)
	Filter(ctx context.Context, req *model.TodoFilter) (*model.TodoFilterResult, error)
	Update(ctx context.Context, id, updaterID uuid.UUID, req model.TodoUpdateRequest) (*model.TodoUpdateResponse, error)
	Delete(ctx context.Context, id, updaterID uuid.UUID) error
	GetOneFlexible(ctx context.Context, field string, value interface{}) (*model.Todo, error)
}

func (s *TodoService) Create(ctx context.Context, req model.CreateTodoRequest, adminID uuid.UUID) (*model.CreateTodoResponse, error) {
	txWithTimeout, cancel := s.todoRepo.DBWithTimeout(ctx)
	defer cancel()

	tx := txWithTimeout.Begin()
	defer tx.Rollback()

	ob := mapper.MapTodo(req, adminID)

	todo, err := s.todoRepo.Create(tx, ob)
	if err != nil {
		return nil, err
	}

	result := model.CreateTodoResponse{
		Meta: model.NewMetaData(),
		Data: model.CreateTodoDataResponse{
			CreatorID: *todo.CreatorID,
			Name:      todo.Name,
			Key:       todo.Key,
			IsActive:  todo.IsActive,
			Code:      todo.Code,
		},
	}

	tx.Commit()
	return &result, nil
}

func (s *TodoService) Get(ctx context.Context, id uuid.UUID) (*model.Todo, error) {
	tx, cancel := s.todoRepo.DBWithTimeout(ctx)
	defer cancel()

	return s.todoRepo.Get(tx, id)
}

func (s *TodoService) Filter(ctx context.Context, req *model.TodoFilter) (*model.TodoFilterResult, error) {
	tx, cancel := s.todoRepo.DBWithTimeout(ctx)
	defer cancel()

	return s.todoRepo.Filter(tx, req)
}

func (s *TodoService) Update(ctx context.Context, id, updaterID uuid.UUID, req model.TodoUpdateRequest) (*model.TodoUpdateResponse, error) {
	txWithTimeout, cancel := s.todoRepo.DBWithTimeout(ctx)
	defer cancel()

	tx := txWithTimeout.Begin()
	defer tx.Rollback()

	ob, err := s.todoRepo.Update(tx, id, updaterID, &req)
	if err != nil {
		return nil, err
	}

	result := model.TodoUpdateResponse{
		Meta: model.NewMetaData(),
		Data: model.TodoUpdateDataResponse{
			Name:        &ob.Name,
			Key:         &ob.Key,
			IsActive:    &ob.IsActive,
			Code:        &ob.Code,
			Description: &ob.Description,
		},
	}

	tx.Commit()
	return &result, nil
}

func (s *TodoService) Delete(ctx context.Context, id, updaterID uuid.UUID) error {
	txWithTimeout, cancel := s.todoRepo.DBWithTimeout(ctx)
	defer cancel()

	tx := txWithTimeout.Begin()
	defer tx.Rollback()

	err := s.todoRepo.Delete(tx, id, updaterID)
	if err != nil {
		return err
	}

	tx.Commit()
	return nil
}

func (s *TodoService) GetOneFlexible(ctx context.Context, field string, value interface{}) (*model.Todo, error) {
	tx, cancel := s.todoRepo.DBWithTimeout(ctx)
	defer cancel()

	return s.todoRepo.GetOneFlexible(tx, field, value)
}
