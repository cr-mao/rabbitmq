package rabbitmq

import (
	"context"
	pool "github.com/jolestar/go-commons-pool/v2"
)

// conn Pool
var connPoolMap map[string]*pool.ObjectPool

func init() {
	connPoolMap = make(map[string]*pool.ObjectPool)
}

func Init(bindDsn map[string]string, maxConNum int) {
	ctx := context.Background()
	for bind, dsn := range bindDsn {
		c := pool.NewDefaultPoolConfig()
		c.LIFO = false
		c.TestOnBorrow = true
		c.TestOnReturn = true
		c.MaxTotal = maxConNum
		p := pool.NewObjectPool(ctx, &connFactory{dsn: dsn}, c)
		connPoolMap[bind] = p
	}
}
