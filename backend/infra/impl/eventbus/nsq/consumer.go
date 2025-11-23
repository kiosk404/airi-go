package nsq

import (
	"context"
	"fmt"

	"github.com/kiosk404/airi-go/backend/infra/contract/eventbus"
	"github.com/kiosk404/airi-go/backend/pkg/lang/conv"
	"github.com/kiosk404/airi-go/backend/pkg/lang/signal"
	"github.com/kiosk404/airi-go/backend/pkg/logs"
	"github.com/kiosk404/airi-go/backend/pkg/utils/safego"
	"github.com/nsqio/go-nsq"
)

func RegisterConsumer(nameServer, topic, group string, consumerHandler eventbus.ConsumerHandler, opts ...eventbus.ConsumerOpt) error {
	if nameServer == "" {
		return fmt.Errorf("name server is empty")
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

	config := nsq.NewConfig()

	consumer, err := nsq.NewConsumer(topic, group, config)
	if err != nil {
		return fmt.Errorf("create consumer failed, err=%w", err)
	}

	consumer.AddHandler(&MessageHandler{
		Topic:           topic,
		Group:           group,
		ConsumerHandler: consumerHandler,
	})

	if err := consumer.ConnectToNSQD(nameServer); err != nil {
		return fmt.Errorf("connect to nsqd failed, err=%w", err)
	}

	safego.Go(context.Background(), func() {
		signal.WaitExit()
		consumer.Stop()
	})

	return nil
}

// Customize the Handler to handle each message received
type MessageHandler struct {
	Topic           string
	Group           string
	ConsumerHandler eventbus.ConsumerHandler
}

func (h *MessageHandler) HandleMessage(m *nsq.Message) error {
	msg := &eventbus.Message{
		Topic: h.Topic,
		Group: h.Group,
		Body:  m.Body,
	}

	logs.Debug("[Subscribe] receive msg : %v \n", conv.DebugJsonToStr(msg))
	err := h.ConsumerHandler.HandleMessage(context.Background(), msg)
	if err != nil {
		logs.Error("[Subscribe] handle msg failed, topic : %s , group : %s, err: %v \n", msg.Topic, msg.Group, err)
		return err
	}

	logs.Debug("subscribe callback: %v \n", conv.DebugJsonToStr(msg))

	return nil
}
