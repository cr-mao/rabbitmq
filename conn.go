package rabbitmq

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

type MqConn struct {
	*amqp.Connection
	*amqp.Channel
}

func NewMqConn(url string) (*MqConn, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	channel, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	go func() {
		<-channel.NotifyClose(make(chan *amqp.Error))
		conn.Close()
	}()

	return &MqConn{Connection: conn, Channel: channel}, nil
}
