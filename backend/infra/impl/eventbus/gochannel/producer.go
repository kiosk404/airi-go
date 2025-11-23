package gochannel

import (
	"context"
	"fmt"
	"sync"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	"github.com/kiosk404/airi-go/backend/infra/contract/eventbus"
)

var (
	pubSub    *gochannel.GoChannel
	initOnce  sync.Once
	startOnce sync.Once
	logger    watermill.LoggerAdapter
	router    *message.Router
)

// producerImpl 生产者实现
type producerImpl struct {
	nameServer string
	topic      string
	publisher  message.Publisher
}

// InitGoChannel 初始化 GoChannel
func InitGoChannel() error {
	var err error
	initOnce.Do(func() {
		logger = watermill.NewStdLogger(false, false)
		pubSub = gochannel.NewGoChannel(gochannel.Config{
			Persistent:                     true,
			BlockPublishUntilSubscriberAck: false,
			OutputChannelBuffer:            0,
		}, logger)
		router, err = message.NewRouter(message.RouterConfig{}, logger)
		if err != nil {
			err = fmt.Errorf("create router failed: %w", err)
			return
		}
	})
	return err
}

// Start 启动
func Start() error {
	var err error
	startOnce.Do(func() {
		if router == nil {
			err = fmt.Errorf("router not initialized")
		}

		go func() {
			if routerErr := router.Run(context.Background()); routerErr != nil {
				logger.Error("Router error", routerErr, nil)
				err = routerErr
			}
		}()

		<-router.Running()

	})
	return err
}

// Close 关闭（应用退出时调用）
func Close() error {
	if router != nil {
		if err := router.Close(); err != nil {
			return err
		}
	}
	if pubSub != nil {
		return pubSub.Close()
	}
	return nil
}

func NewProducer(nameServer, topic, group string) (eventbus.Producer, error) {
	// GoChannel 不需要 nameServer 和 group，但保留参数以兼容接口

	Start()

	if topic == "" {
		return nil, fmt.Errorf("topic is empty")
	}

	if pubSub == nil {
		return nil, fmt.Errorf("not initialized, call Init() first")
	}

	return &producerImpl{
		nameServer: nameServer,
		topic:      topic,
		publisher:  pubSub,
	}, nil
}

// Send 发送单条消息
func (p *producerImpl) Send(ctx context.Context, body []byte, opts ...eventbus.ProduceOpt) error {
	msg := message.NewMessage(watermill.NewUUID(), body)
	msg.SetContext(ctx)

	err := p.publisher.Publish(p.topic, msg)
	if err != nil {
		return fmt.Errorf("[producerImpl] send message failed: %w", err)
	}
	return nil
}

// BatchSend 批量发送消息
func (p *producerImpl) BatchSend(ctx context.Context, bodyArr [][]byte, opts ...eventbus.ProduceOpt) error {
	option := eventbus.ProduceOption{}
	for _, opt := range opts {
		opt(&option)
	}

	messages := make([]*message.Message, len(bodyArr))
	for i, body := range bodyArr {
		msg := message.NewMessage(watermill.NewUUID(), body)
		msg.SetContext(ctx)
		messages[i] = msg
	}

	err := p.publisher.Publish(p.topic, messages...)
	if err != nil {
		return fmt.Errorf("[BatchSend] send message failed: %w", err)
	}
	return nil
}
