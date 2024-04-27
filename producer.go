package rabbitmq

import (
	"context"

	amqp "github.com/rabbitmq/amqp091-go"
)

// 发布
func Pub(ctx context.Context, bind string, exchange string, routingKey string, message []byte) error {
	sess, returnFn, err := Borrow(ctx, bind)
	if err != nil {
		return err
	}
	defer returnFn(ctx)
	err = sess.Channel.PublishWithContext(ctx,
		exchange,
		routingKey,
		false, // mandatory 如果为 true,在 exchange 正常且可到达的情况，如果 exchange+routekey 无法投递给 queue，那么 mq 会将消息返还给生产者.
		false, // immediate
		amqp.Publishing{
			ContentType:  "text/plain",
			DeliveryMode: amqp.Persistent, //持久化，
			Body:         message,
		},
	)
	return err
}
