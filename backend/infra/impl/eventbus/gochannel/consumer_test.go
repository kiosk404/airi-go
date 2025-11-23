package gochannel

import (
	"context"
	"fmt"
	"testing"

	"github.com/kiosk404/airi-go/backend/infra/contract/eventbus"
)

type UserCreatedHandler struct{}

func (h *UserCreatedHandler) HandleMessage(ctx context.Context, msg *eventbus.Message) error {
	fmt.Printf("Received: topic=%s, group=%s, body=%s\n",
		msg.Topic, msg.Group, string(msg.Body))
	return nil
}

func TestRegisterConsumer(t *testing.T) {
	// 1. 初始化
	if err := InitGoChannel(); err != nil {
		panic(err)
	}
	defer Close()
	handler := &UserCreatedHandler{}
	if err := RegisterConsumer("", "user.created", "notification", handler); err != nil {
		panic(err)
	}
	if err := RegisterConsumer("", "user.updated", "notification", handler); err != nil {
		panic(err)
	}
	producer1, err := NewProducer("", "user.created", "notification")
	if err != nil {
		panic(err)
	}
	producer2, err := NewProducer("", "user.updated", "notification")
	if err != nil {
		panic(err)
	}

	if err := producer1.Send(context.Background(), []byte(`{"user_id":"123"}`)); err != nil {
		panic(err)
	}

	if err := producer1.Send(context.Background(), []byte(`{"user_id":"678"}`)); err != nil {
		panic(err)
	}

	messages := [][]byte{
		[]byte(`{"user_id":"123"}`),
		[]byte(`{"user_id":"456"}`),
		[]byte(`{"user_id":"789"}`),
	}
	if err := producer2.BatchSend(context.Background(), messages); err != nil {
		panic(err)
	}

	select {}

}
