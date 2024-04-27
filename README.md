## golang rabbitmq client 

golang 操作 rabbitmq封装

- 带有连接池实现（单连接在高并发项目下，是不够用的)
- 统一使用topic模式，通用能覆盖所有场景。 

```shell
go get github.com/cr-mao/rabbitmq
```

### example

消费者
```go
package main 
import (
    "github.com/cr-mao/rabbitmq"
)
const bindClient string ="default"

func main() {
	ctx := context.Background()
    // 
	var bindDsn = map[string]string{
		bindClient: "amqp://admin:admin@127.0.0.1:5672/",
	}
	rabbitmq.Init(bindDsn, 10)
	var consumer = &rabbitmq.AMQPConsumer{
		Bind: bindClient,
		Subscription: &rabbitmq.Subscription{
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
```

生产者:
```go
package main
import (
	"github.com/cr-mao/rabbitmq"
)

const bindClient string ="default"

func main(){
	var bindDsn = map[string]string{
		bindClient: "amqp://admin:admin@127.0.0.1:5672/",
	}
	rabbitmq.Init(bindDsn, 10)
	ctx := context.Background()
	rabbitmq.Pub(ctx, bindClient, "testexchange", "test", []byte("cr-mao"))
}
```

### links 

[golang通用连接池](github.com/jolestar/go-commons-pool/v2)