package main

import (
	"context"
	"fmt"
	"rabbitmq-demo/rabbitmq"
	"time"
)

func main() {

	//ctx := context.TODO()
	//producerExchange(ctx)
	//producerQueue(ctx)

	select {}
}

func producerExchange(ctx context.Context) {
	//producerEx
	producerEx, err := rabbitmq.NewProducer(ctx, rabbitmq.URLConfig{
		UserName: "guest",
		PassWord: "guest",
		IPAddr:   "127.0.0.1",
		Port:     5672,
		Vhost:    "",
	}, rabbitmq.AMQP_MODE_EXCHANGE, "logs", "")
	if err != nil {
		panic(err)
	}

	go func() {
		if _, err := rabbitmq.NewConsumer(ctx, rabbitmq.URLConfig{
			IPAddr:   "127.0.0.1",
			Port:     5672,
			Vhost:    "",
			UserName: "guest",
			PassWord: "guest",
		}, rabbitmq.AMQP_MODE_EXCHANGE, "logs", "", func(data []byte) {
			fmt.Println("a revData", string(data))
		}); err != nil {
			fmt.Println(err)
		}
	}()
	// name 1
	go func() {
		if _, err := rabbitmq.NewConsumer(ctx, rabbitmq.URLConfig{
			IPAddr:   "127.0.0.1",
			Port:     5672,
			Vhost:    "",
			UserName: "guest",
			PassWord: "guest",
		}, rabbitmq.AMQP_MODE_EXCHANGE, "logs", "", func(data []byte) {
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
}

func producerQueue(ctx context.Context) {
	producerQueue, err := rabbitmq.NewProducer(ctx, rabbitmq.URLConfig{
		UserName: "guest",
		PassWord: "guest",
		IPAddr:   "127.0.0.1",
		Port:     5672,
		Vhost:    "",
	}, rabbitmq.AMQP_MODE_QUEUE, "", "queue_logs")
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
		if _, err := rabbitmq.NewConsumer(ctx, rabbitmq.URLConfig{
			IPAddr:   "127.0.0.1",
			Port:     5672,
			Vhost:    "",
			UserName: "guest",
			PassWord: "guest",
		}, rabbitmq.AMQP_MODE_QUEUE, "", "queue_logs", func(data []byte) {
			fmt.Println("queue revData", string(data))
		}); err != nil {
			fmt.Println(err)
		}
	}()
}
