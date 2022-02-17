package consumer

import (
	"log"
	"nats-workers/models"
	"time"

	cbor "github.com/fxamacker/cbor/v2"
	"github.com/thiagozs/go-nats/nats"
	"github.com/thiagozs/go-nats/streaming"
)

type WorkerConsumerRepo interface {
	GetMessage() <-chan []byte
	Consumer()
	ExitConsumer()
}

type WorkerConsumer struct {
	stream     streaming.StreamingServiceRepo
	channel    string
	workerStop chan struct{}
}

func New(instance nats.NatsServiceRepo, channel string) WorkerConsumerRepo {

	stm := streaming.NewService(instance)

	if err := stm.Subscribe(channel); err != nil {
		log.Fatalf("consumer subscribe: %+v\n", err)
	}

	return &WorkerConsumer{
		stream:     stm,
		channel:    channel,
		workerStop: make(chan struct{}),
	}
}

func (w *WorkerConsumer) GetMessage() <-chan []byte {
	return w.stream.GetMessage(w.channel)
}

func (w *WorkerConsumer) ExitConsumer() {
	w.workerStop <- struct{}{}
}

func (w *WorkerConsumer) Consumer() {
	for {
		select {
		case msg := <-w.GetMessage():
			var message models.Message
			if err := cbor.Unmarshal(msg, &message); err != nil {
				log.Fatalf("consumer unmarshal cbor: %+v\n", err)
			}
			log.Printf("%+v\n", message)

		case <-time.After(time.Second * 1):
			log.Println("timeout")

		case <-w.workerStop:
			log.Println("exit consumer")
			return
		}
	}
}
