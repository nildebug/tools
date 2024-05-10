package rabbitmq

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestNewProducer(t *testing.T) {
	ctx := context.Background()
	producerEx, err := NewProducer(ctx, URLConfig{
		UserName: "guest",
		PassWord: "guest",
		IPAddr:   "127.0.0.1",
		Port:     5672,
		Vhost:    "",
	}, AMQP_MODE_EXCHANGE, "logs", "")
	if err != nil {
		panic(err)
	}

	go func() {
		if _, err := NewConsumer(ctx, URLConfig{
			IPAddr:   "127.0.0.1",
			Port:     5672,
			Vhost:    "",
			UserName: "guest",
			PassWord: "guest",
		}, AMQP_MODE_EXCHANGE, "logs", "", func(data []byte) {
			fmt.Println("a revData", string(data))
		}); err != nil {
			fmt.Println(err)
		}
	}()
	// name 1
	go func() {
		if _, err := NewConsumer(ctx, URLConfig{
			IPAddr:   "127.0.0.1",
			Port:     5672,
			Vhost:    "",
			UserName: "guest",
			PassWord: "guest",
		}, AMQP_MODE_EXCHANGE, "logs", "", func(data []byte) {
			fmt.Println("b revData", string(data))
		}); err != nil {
			fmt.Println(err)
		}
	}()
	// name 2
	go func() {
		for i := 0; i < 100; i++ {
			if err := producerEx.SendMessage(ctx, []byte(fmt.Sprintf("hello %d", i))); err != nil {
				fmt.Println("sendMessage fail", err)
			}
			time.Sleep(1 * time.Second)
		}
	}()

	select {}
}

func TestProducerQueue(t *testing.T) {
	ctx := context.Background()
	producerQueue, err := NewProducer(ctx, URLConfig{
		UserName: "guest",
		PassWord: "guest",
		IPAddr:   "127.0.0.1",
		Port:     5672,
		Vhost:    "",
	}, AMQP_MODE_QUEUE, "", "queue_logs")
	if err != nil {
		panic(err)
	}

	go func() {
		for i := 0; i < 100; i++ {
			if err := producerQueue.SendMessage(ctx, []byte(fmt.Sprintf("hello %d", i))); err != nil {
				fmt.Println("sendMessage fail", err)
			}
			time.Sleep(1 * time.Second)
		}
	}()

	go func() {
		if _, err := NewConsumer(ctx, URLConfig{
			IPAddr:   "127.0.0.1",
			Port:     5672,
			Vhost:    "",
			UserName: "guest",
			PassWord: "guest",
		}, AMQP_MODE_QUEUE, "", "queue_logs", func(data []byte) {
			fmt.Println("queue revData", string(data))
		}); err != nil {
			fmt.Println(err)
		}
	}()

	go func() {
		if _, err := NewConsumer(ctx, URLConfig{
			IPAddr:   "127.0.0.1",
			Port:     5672,
			Vhost:    "",
			UserName: "guest",
			PassWord: "guest",
		}, AMQP_MODE_QUEUE, "", "queue_logs", func(data []byte) {
			fmt.Println("b queue revData", string(data))
		}); err != nil {
			fmt.Println(err)
		}
	}()

	select {}
}
