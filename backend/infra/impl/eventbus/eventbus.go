package eventbus

import (
	"fmt"
	"os"

	"github.com/kiosk404/airi-go/backend/infra/contract/eventbus"
	"github.com/kiosk404/airi-go/backend/infra/impl/eventbus/gochannel"
	"github.com/kiosk404/airi-go/backend/infra/impl/eventbus/nsq"
	"github.com/kiosk404/airi-go/backend/infra/impl/eventbus/rmq"
	"github.com/kiosk404/airi-go/backend/types/consts"
)

type (
	Producer        = eventbus.Producer
	ConsumerService = eventbus.ConsumerService
	ConsumerHandler = eventbus.ConsumerHandler
	ConsumerOpt     = eventbus.ConsumerOpt
	Message         = eventbus.Message
)

type consumerServiceImpl struct{}

func NewConsumerService() ConsumerService {
	return &consumerServiceImpl{}
}

func DefaultSVC() ConsumerService {
	return eventbus.GetDefaultSVC()
}

func (c consumerServiceImpl) RegisterConsumer(nameServer, topic, group string, consumerHandler eventbus.ConsumerHandler, opts ...eventbus.ConsumerOpt) error {
	tp := os.Getenv(consts.MQTypeKey)
	switch tp {
	case "nsq":
		return nsq.RegisterConsumer(nameServer, topic, group, consumerHandler, opts...)
	case "rmq":
		return rmq.RegisterConsumer(nameServer, topic, group, consumerHandler, opts...)
	case "gochannel":
		if err := gochannel.InitGoChannel(); err != nil {
			return err
		}
		return gochannel.RegisterConsumer(nameServer, topic, group, consumerHandler, opts...)
	}

	return fmt.Errorf("invalid mq type: %s , only support nsq, kafka, rmq", tp)
}

func NewProducer(nameServer, topic, group string, retries int) (eventbus.Producer, error) {
	tp := os.Getenv(consts.MQTypeKey)
	switch tp {
	case "nsq":
		return nsq.NewProducer(nameServer, topic, group)
	case "rmq":
		return rmq.NewProducer(nameServer, topic, group, retries)
	case "gochannel":
		return gochannel.NewProducer(nameServer, topic, group)
	}

	return nil, fmt.Errorf("invalid mq type: %s , only support nsq, kafka, rmq", tp)
}

func InitResourceEventBusProducer() (eventbus.Producer, error) {
	nameServer := os.Getenv(consts.MQServer)
	resourceEventBusProducer, err := NewProducer(nameServer,
		consts.RMQTopicResource, consts.RMQConsumeGroupResource, 1)
	if err != nil {
		return nil, fmt.Errorf("init resource producer failed, err=%w", err)
	}

	return resourceEventBusProducer, nil
}

func InitAppEventProducer() (eventbus.Producer, error) {
	nameServer := os.Getenv(consts.MQServer)
	appEventProducer, err := NewProducer(nameServer, consts.RMQTopicApp, consts.RMQConsumeGroupApp, 1)
	if err != nil {
		return nil, fmt.Errorf("init app producer failed, err=%w", err)
	}

	return appEventProducer, nil
}

func InitKnowledgeEventBusProducer() (eventbus.Producer, error) {
	nameServer := os.Getenv(consts.MQServer)

	knowledgeProducer, err := NewProducer(nameServer, consts.RMQTopicKnowledge, consts.RMQConsumeGroupKnowledge, 2)
	if err != nil {
		return nil, fmt.Errorf("init knowledge producer failed, err=%w", err)
	}

	return knowledgeProducer, nil
}
