# Nats Wrapper with consumer and publisher

Simple implementation of NATS streaming client with consumer and publisher.

```go

// Loadding server configs
// follow the main.go to see how to use

	c := consumer.New(ns, "test")
	p := publisher.New(ns, "test")

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

```

## Versioning and license

We use SemVer for versioning. You can see the versions available by checking the tags on this repository.

For more details about our license model, please take a look at the LICENSE file

---

2022, thiagozs