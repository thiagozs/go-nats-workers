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

	defer c.ExitConsumer()

	go func() {
		c.Consumer()
	}()

	for i := 1; i <= 11; i++ {
		m := models.Message{
			Order:       fmt.Sprintf("%s-%d", "order", i),
			Matches:     fmt.Sprintf("%s-%d", "matches", i),
			Operations:  fmt.Sprintf("%s-%d", "operations", i),
			Commissions: fmt.Sprintf("%s-%d", "commissions", i),
		}
		if err := p.SentMessage(m); err != nil {
			log.Fatal(err)
		}
	}
}
