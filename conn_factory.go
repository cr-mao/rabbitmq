// 连接工程，创建连接、销毁连接、验证连接
package rabbitmq

import (
	"context"

	pool "github.com/jolestar/go-commons-pool/v2"
)

// connFactory new conn from dsn, and wrap it to pool.PooledObject
type connFactory struct {
	dsn string
}

func (sf *connFactory) MakeObject(ctx context.Context) (*pool.PooledObject, error) {
	conn, err := NewMqConn(sf.dsn)
	return pool.NewPooledObject(conn), err
}

func (sf *connFactory) DestroyObject(ctx context.Context, object *pool.PooledObject) error {
	object.Object.(*MqConn).Connection.Close()
	return nil
}

func (sf *connFactory) ValidateObject(ctx context.Context, object *pool.PooledObject) bool {
	return !object.Object.(*MqConn).Connection.IsClosed()
}

func (sf *connFactory) ActivateObject(ctx context.Context, object *pool.PooledObject) error {
	return nil
}

func (sf *connFactory) PassivateObject(ctx context.Context, object *pool.PooledObject) error {
	return nil
}
