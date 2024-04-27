/**
* @Author: maozhongyu
* @Desc:  测试用例
* @Date: 2024/4/27
**/
package rabbitmq

import (
	"context"
	"fmt"
	"testing"
)

const (
	bindClient string = "default"
)

// admin,admin 账号密码是新建的
func TestProduct(t *testing.T) {
	var bindDsn = map[string]string{
		bindClient: "amqp://admin:admin@127.0.0.1:5672/",
	}
	Init(bindDsn, 10)
	ctx := context.Background()

	Pub(ctx, bindClient, "testexchange", "test", []byte("cr-mao"))
}

func TestConsumer(t *testing.T) {
	ctx := context.Background()
	var bindDsn = map[string]string{
		bindClient: "amqp://admin:admin@127.0.0.1:5672/",
	}
	Init(bindDsn, 10)
	var consumer = &AMQPConsumer{
		Bind: bindClient,
		Subscription: &Subscription{
			Exchange: "testexchange",
			Queue:    "testqueue",
			Key: []string{
				"test",
			},
		},
		MessageHandler: func(ctx context.Context, msg []byte) error {
			fmt.Println(string(msg))
			return nil
		},
	}
	err := consumer.Run(ctx)
	if err != nil {
		fmt.Println(err)
	}
}

//func Test
