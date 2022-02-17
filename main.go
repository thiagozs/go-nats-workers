package main

import (
	"fmt"
	"log"
	"nats-workers/models"
	"nats-workers/workers/consumer"
	"nats-workers/workers/publisher"
	"time"

	"github.com/thiagozs/go-nats/nats"
)

func main() {
	loadingServer := func() (nats.NatsServiceRepo, error) {
		opts := []nats.Options{
			nats.ClusterName("test-cluster"),
			nats.ClusterID("clusterID"),
			nats.NastPort("4223"),
			nats.NastURL("localhost"),
			nats.PubAckWait(30 * time.Second),
			nats.MaxInFlight(20),
			nats.PingParams(nats.PingParameter{
				TTLA: 100,
				TTLB: 250,
			}),
		}
		ns, err := nats.NewService(opts...)
		if err != nil {
			return nil, err
		}
		return ns, nil
	}

	ns, err := loadingServer()
	if err != nil {
		log.Fatal(err)
	}

	c := consumer.New(ns, "test")
	p := publisher.New(ns, "test")

	go func() {
		c.Consumer()
	}()

	for i := 0; i <= 10; i++ {
		incIndex := i + 1
		m := models.Message{
			Order:       fmt.Sprintf("%s-%d", "order", incIndex),
			Matches:     fmt.Sprintf("%s-%d", "matches", incIndex),
			Operations:  fmt.Sprintf("%s-%d", "operations", incIndex),
			Commissions: fmt.Sprintf("%s-%d", "commissions", incIndex),
		}
		if err := p.SentMessage(m); err != nil {
			log.Fatal(err)
		}
	}
	c.ExitConsumer()
}
