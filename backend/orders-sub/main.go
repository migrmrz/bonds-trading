package main

import (
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
)

const ordersSubject = "com.trading.orders.update"

func main() {
	// subscribe to NATS subject
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		panic(err)
	}

	defer nc.Close()

	nc.Subscribe(ordersSubject, func(m *nats.Msg) {
		data := m.Data
		fmt.Println(string(data))
	})

	// leave service running for an hour
	time.Sleep(1 * time.Hour)
}
