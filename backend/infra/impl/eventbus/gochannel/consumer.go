package gochannel

import (
	"fmt"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/kiosk404/airi-go/backend/infra/contract/eventbus"
)

func RegisterConsumer(nameServer, topic, group string, consumerHandler eventbus.ConsumerHandler, opts ...eventbus.ConsumerOpt) error {
	if router == nil {
		return fmt.Errorf("not initialized, call InitGoChannel first")
	}

	if topic == "" {
		return fmt.Errorf("topic is empty")
	}

	if group == "" {
		return fmt.Errorf("group is empty")
	}

	if consumerHandler == nil {
		return fmt.Errorf("consumer handler is nil")
	}

	// 处理选项（当前为空，预留扩展）
	options := &eventbus.ConsumerOption{}
	for _, opt := range opts {
		opt(options)
	}

	// 创建处理函数
	handlerFunc := func(msg *message.Message) error {
		eventMsg := &eventbus.Message{
			Topic: topic,
			Group: group,
			Body:  msg.Payload,
		}

		err := consumerHandler.HandleMessage(msg.Context(), eventMsg)
		if err != nil {
			logger.Error("Handle message failed", err, watermill.LogFields{
				"topic": topic,
				"group": group,
			})
			return err
		}

		logger.Debug("Message handled successfully", watermill.LogFields{
			"topic": topic,
			"group": group,
		})

		msg.Ack()
		return nil
	}

	// 注册到路由器
	handlerName := fmt.Sprintf("%s:%s", topic, group)
	router.AddConsumerHandler(
		handlerName,
		topic,
		pubSub,
		handlerFunc,
	)

	return nil
}

// MessageHandler Customize the Handler to handle each message received
type MessageHandler struct {
	Topic           string
	Group           string
	ConsumerHandler eventbus.ConsumerHandler
}
