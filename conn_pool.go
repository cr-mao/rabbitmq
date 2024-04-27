package rabbitmq

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
)

// 拿连接，闭包归还准备
func Borrow(ctx context.Context, bind string) (*MqConn, func(context.Context), error) {
	p, ok := connPoolMap[bind]
	if !ok {
		return nil, nil, errors.New(fmt.Sprintf("bind not found:%s", bind))
	}
	conn, err := p.BorrowObject(ctx)
	if err != nil {
		return nil, nil, errors.Wrap(err, "Borrow RabbitMQ MqConn from pool failed.")
	}

	return conn.(*MqConn), func(ctx context.Context) {
		p.ReturnObject(ctx, conn)
	}, nil
}
