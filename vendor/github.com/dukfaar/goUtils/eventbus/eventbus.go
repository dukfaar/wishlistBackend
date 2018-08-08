package eventbus

import (
	"encoding/json"

	"github.com/nsqio/go-nsq"
)

type EventHandler func(payload []byte) error

type EventBus interface {
	Emit(topic string, payload interface{}) error
	On(topic string, channel string, handler EventHandler) Consumer
}

type Consumer interface {
	Stop()
}

type NsqEventBus struct {
	config        *nsq.Config
	producer      *nsq.Producer
	nsqdUrl       string
	nsqLookupdUrl string
}

type NsqConsumer struct {
	consumer *nsq.Consumer
}

func (c *NsqConsumer) Stop() {
	c.consumer.Stop()
}

func NewNsqEventBus(nsqdUrl string, nsqLookupdUrl string) *NsqEventBus {
	config := nsq.NewConfig()
	producer, err := nsq.NewProducer(nsqdUrl, config)

	if err != nil {
		panic(err)
	}

	return &NsqEventBus{
		config:        config,
		nsqLookupdUrl: nsqLookupdUrl,
		nsqdUrl:       nsqdUrl,
		producer:      producer,
	}
}

func (b *NsqEventBus) Emit(topic string, payload interface{}) error {
	message, err := json.Marshal(&payload)

	if err != nil {
		return err
	}

	return b.producer.Publish(topic, message)
}

func (b *NsqEventBus) On(topic string, channel string, handler EventHandler) Consumer {
	consumer, err := nsq.NewConsumer(topic, channel, b.config)

	if err != nil {
		return nil
	}

	consumer.AddHandler(nsq.HandlerFunc(func(message *nsq.Message) error {
		go handler(message.Body)
		return nil
	}))

	if err := consumer.ConnectToNSQLookupd(b.nsqLookupdUrl); err != nil {
		return nil
	}

	return &NsqConsumer{
		consumer: consumer,
	}
}
