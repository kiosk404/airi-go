package eventbus

import (
	"context"
)

//go:generate  mockgen -destination ./mock/eventbus_mock.go -package mock -source eventbus.go Factory
type Producer interface {
	Send(ctx context.Context, body []byte, opts ...ProduceOpt) error
	BatchSend(ctx context.Context, bodyArr [][]byte, opts ...ProduceOpt) error
}

var defaultSVC ConsumerService

func SetDefaultSVC(svc ConsumerService) {
	defaultSVC = svc
}

func GetDefaultSVC() ConsumerService {
	return defaultSVC
}

type ConsumerService interface {
	RegisterConsumer(nameServer, topic, group string, consumerHandler ConsumerHandler, opts ...ConsumerOpt) error
}

type ConsumerHandler interface {
	HandleMessage(ctx context.Context, msg *Message) error
}

type Message struct {
	Topic string
	Group string
	Body  []byte
}
