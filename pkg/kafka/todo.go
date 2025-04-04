package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
	"good-template-go/pkg/model"
	"good-template-go/pkg/repo"
)

type TodoKafkaHandlers struct {
	repoTodo repo.RepoTodoInterface
}

func NewTodoKafkaHandlers(
	repoTodo repo.RepoTodoInterface) *TodoKafkaHandlers {
	return &TodoKafkaHandlers{
		repoTodo: repoTodo,
	}
}

func (k *TodoKafkaHandlers) KafkaProcess(ctx context.Context, kafkaMessage kafka.Message) error {
	var message model.TodoKafkaMessage
	err := json.Unmarshal(kafkaMessage.Value, &message)
	if err != nil {
		return fmt.Errorf("error unmarshalling Kafka message: %v", err)
	}

	oldData := message.Payload.Before
	newData := message.Payload.After

	switch {
	// delete
	case oldData != nil && newData == nil:
		// ignore
		return nil

	// update
	case oldData != nil && newData != nil:
		err = k.repoTodo.Update(ctx, newData.ID, newData)
		if err != nil {
			return err
		}
	// create
	case oldData == nil && newData != nil:
		err = k.repoTodo.Create(ctx, newData)
		if err != nil {
			return err
		}
	default:
		// ignore
		return nil
	}

	return nil
}
