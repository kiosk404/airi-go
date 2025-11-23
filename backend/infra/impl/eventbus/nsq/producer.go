package nsq

import (
	"context"
	"fmt"

	"github.com/kiosk404/airi-go/backend/infra/contract/eventbus"
	"github.com/kiosk404/airi-go/backend/pkg/lang/signal"
	"github.com/kiosk404/airi-go/backend/pkg/utils/safego"
	"github.com/nsqio/go-nsq"
)

type producerImpl struct {
	nameServer string
	topic      string
	p          *nsq.Producer
}

func NewProducer(nameServer, topic, group string) (eventbus.Producer, error) {
	if nameServer == "" {
		return nil, fmt.Errorf("name server is empty")
	}

	if topic == "" {
		return nil, fmt.Errorf("topic is empty")
	}

	config := nsq.NewConfig()

	producer, err := nsq.NewProducer(nameServer, config)
	if err != nil {
		return nil, fmt.Errorf("create producer failed, err=%w", err)
	}

	safego.Go(context.Background(), func() {
		signal.WaitExit()
		producer.Stop()
	})

	return &producerImpl{
		nameServer: nameServer,
		topic:      topic,
		p:          producer,
	}, nil
}

func (r *producerImpl) Send(ctx context.Context, body []byte, opts ...eventbus.ProduceOpt) error {
	err := r.p.Publish(r.topic, body)
	if err != nil {
		return fmt.Errorf("[producerImpl] send message failed: %w", err)
	}
	return err
}

func (r *producerImpl) BatchSend(ctx context.Context, bodyArr [][]byte, opts ...eventbus.ProduceOpt) error {
	option := eventbus.ProduceOption{}
	for _, opt := range opts {
		opt(&option)
	}

	err := r.p.MultiPublish(r.topic, bodyArr)
	if err != nil {
		return fmt.Errorf("[BatchSend] send message failed: %w", err)
	}
	return nil
}
