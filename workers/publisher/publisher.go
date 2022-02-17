package publisher

import (
	"log"

	"github.com/thiagozs/go-nats/nats"
	"github.com/thiagozs/go-nats/streaming"

	cbor "github.com/fxamacker/cbor/v2"
)

type WorkerPublisherRepo interface {
	SentMessage(model interface{}) error
}

type WorkerPublisher struct {
	stream  streaming.StreamingServiceRepo
	channel string
}

func New(instance nats.NatsServiceRepo, channel string) WorkerPublisherRepo {

	stm := streaming.NewService(instance)

	if err := stm.Subscribe(channel); err != nil {
		log.Fatalf("publisher: %+v\n", err)
	}

	return &WorkerPublisher{
		stream:  stm,
		channel: channel,
	}
}

func (w *WorkerPublisher) SentMessage(model interface{}) error {

	payload, err := cbor.Marshal(model)
	if err != nil {
		return err
	}
	sms := nats.Message{
		Payload: payload,
		Subject: w.channel,
	}
	return w.stream.Publish(sms)
}
