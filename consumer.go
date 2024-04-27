package rabbitmq

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	amqp "github.com/rabbitmq/amqp091-go"
)

// 函数带了上下文，可以用来中间件获取一些内容，如执行时间。
type MessageHandler func(ctx context.Context, msg []byte) error

type Subscription struct {
	Exchange string
	Queue    string
	Key      []string
}

func (m *Subscription) GetExchange() string {
	if m != nil {
		return m.Exchange
	}
	return ""
}

func (m *Subscription) GetQueue() string {
	if m != nil {
		return m.Queue
	}
	return ""
}

func (m *Subscription) GetKey() []string {
	if m != nil {
		return m.Key
	}
	return nil
}

func AMQPTopology(channel *amqp.Channel, sub *Subscription) error {
	if err := channel.ExchangeDeclare(
		sub.GetExchange(), // name of the exchange
		"topic",           // type
		true,              // durable
		false,             // delete when complete
		false,             // internal
		false,             // noWait
		nil,               // arguments
	); err != nil {
		return fmt.Errorf("Exchange Declare: %s", err)
	}
	queue, err := channel.QueueDeclare(
		sub.GetQueue(), // name of the queue
		true,           // durable
		false,          // delete when unused
		false,          // exclusive
		false,          // noWait
		nil,            // arguments
	)
	if err != nil {
		return fmt.Errorf("Queue Declare: %s", err)
	}
	for _, key := range sub.GetKey() {
		if err = channel.QueueBind(
			queue.Name,        // name of the queue
			key,               // bindingKey
			sub.GetExchange(), // sourceExchange
			false,             // noWait
			nil,               // arguments
		); err != nil {
			return fmt.Errorf("Queue Bind: %s", err)
		}
	}
	return nil
}

type AMQPConsumer struct {
	// RabbitMQ Topology
	Bind           string
	Subscription   *Subscription
	MessageHandler MessageHandler
	// Setup Hooks
	SetupHooks []func(*AMQPConsumer, *amqp.Channel) error
}

func (c *AMQPConsumer) Run(ctx context.Context) error {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	done := make(chan error, 1)

	sess, returnFn, err := Borrow(ctx, c.Bind)
	if err != nil {
		return err
	}
	defer returnFn(ctx)
	// Init Hook
	if c.SetupHooks == nil {
		// Default AMQPTopology
		c.SetupHooks = []func(*AMQPConsumer, *amqp.Channel) error{
			func(c *AMQPConsumer, ch *amqp.Channel) error {
				return AMQPTopology(ch, c.Subscription)
			},
		}
	}
	for _, f := range c.SetupHooks {
		err := f(c, sess.Channel)
		if err != nil {
			return err
		}
	}

	deliveries, err := sess.Channel.Consume(
		c.Subscription.GetQueue(), // name
		"",                        // consumerTag,
		false,                     // noAck
		false,                     // exclusive
		false,                     // noLocal
		false,                     // noWait
		nil,                       // arguments
	)
	if err != nil {
		return fmt.Errorf("Queue Consume: %s", err)
	}

	// Stop On Channel Close
	// TODO: Reconnect instead of close
	go func(done chan<- error) {
		<-sess.Channel.NotifyClose(make(chan *amqp.Error))
		done <- nil
	}(done)

	go func(ctx context.Context, deliveries <-chan amqp.Delivery) {
		handler := c.MessageHandler
		for delivery := range deliveries {
			handler(ctx, delivery.Body)
			delivery.Ack(false) // 消息确认 ack。
		}
	}(ctx, deliveries)

	select {
	case <-quit:
	case <-done:
	}
	return nil
}
