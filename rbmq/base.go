package rbmq

import (
	"context"
	"encoding/json"
	"fmt"
	pool "github.com/jolestar/go-commons-pool/v2"
	amqp "github.com/rabbitmq/amqp091-go"
	"sync"
	"time"
)

type (
	Idle struct {
		MaxSize int `json:"max_size" yaml:"max_size"`
		MinIdle int `json:"min_idle" yaml:"min_idle"`
		MaxIdle int `json:"max_idle" yaml:"max_idle"`
	}

	Conf struct {
		Addr     string `json:"addr" yaml:"addr"`
		Port     string `json:"port" yaml:"port"`
		User     string `json:"user" yaml:"user"`
		Pwd      string `json:"pwd" yaml:"pwd"`
		Vhost    string `json:"vhost" yaml:"vhost"`
		PoolIdle Idle   `json:"pool_idle" yaml:"vhost"`
	}

	ChannelFactory struct {
		mqConfig  Conf
		mqConnStr string
	}
)

type DataWithCtx struct {
	Ctx  []CtxData `json:"ctx"`
	Data []byte    `json:"data"`
}

var (
	ChannelPool *pool.ObjectPool // 连接池

	activeConn *Connection // 当前操作链接
	mutex      sync.Mutex  // 互斥锁
)

// Init 初始化
func Init(c Conf) (err error) {
	if c.Addr == "" {
		return fmt.Errorf("RabbitMQ 连接地址为空！")
	}

	connStr := fmt.Sprintf("amqp://%s:%s@%s:%s/%s",
		c.User,
		c.Pwd,
		c.Addr,
		c.Port,
		c.Vhost)

	activeConn, _ = Dial(connStr)

	maxSize := c.PoolIdle.MaxSize
	if maxSize <= 0 || maxSize > 2000 {
		maxSize = 2000
	}
	minIdle := c.PoolIdle.MinIdle
	if minIdle <= 0 || minIdle > 2000 {
		minIdle = 2000
	}
	maxIdle := c.PoolIdle.MaxIdle
	if maxIdle <= 0 || maxIdle > 2000 {
		maxIdle = 2000
	}

	ChannelPool = pool.NewObjectPool(context.TODO(), &ChannelFactory{mqConfig: c, mqConnStr: connStr}, &pool.ObjectPoolConfig{
		LIFO:                     false,
		MaxTotal:                 maxSize,
		MinIdle:                  minIdle,
		MaxIdle:                  maxIdle,
		TestOnCreate:             false,
		TestOnBorrow:             false,
		TestOnReturn:             false,
		TestWhileIdle:            false,
		BlockWhenExhausted:       true,
		MinEvictableIdleTime:     0,
		SoftMinEvictableIdleTime: 0,
		NumTestsPerEvictionRun:   0,
		EvictionPolicyName:       "",
		TimeBetweenEvictionRuns:  0,
		EvictionContext:          nil,
	})

	return nil
}

// ExchangeDeclare 创建交换机.
func (ch *Channel) ExchangeDeclare(name string, kind string) (err error) {
	return ch.Channel.ExchangeDeclare(name, kind, true, false, false, false, nil)
}

// 发布消息. 带上下文传递
func (ch *Channel) Publish(ctx context.Context, exchange, key string, body []byte) (err error) {
	body, err = AddContextToMessage(ctx, body)
	if err != nil {
		return err // 无法添加上下文，直接返回
	}
	_, err = ch.Channel.PublishWithDeferredConfirmWithContext(ctx, exchange, key, false, false,
		amqp.Publishing{ContentType: "text/plain", Body: body})
	return err
}

// Publish 发布消息. 不带上下文传递
func (ch *Channel) PublishWithoutContext(ctx context.Context, exchange, key string, body []byte) (err error) {
	_, err = ch.Channel.PublishWithDeferredConfirmWithContext(ctx, exchange, key, false, false,
		amqp.Publishing{ContentType: "text/plain", Body: body})
	return err
}

// 发布延迟消息. 带上下文传递
func (ch *Channel) PublishDelay(ctx context.Context, exchange, key string, body []byte, timer time.Duration) (err error) {
	body, err = AddContextToMessage(ctx, body)
	if err != nil {
		return err // 无法添加上下文，直接返回
	}
	_, err = ch.Channel.PublishWithDeferredConfirmWithContext(ctx, exchange, key, false, false,
		amqp.Publishing{ContentType: "text/plain", Body: body, Expiration: fmt.Sprintf("%d", timer.Milliseconds())})
	return err
}

// 发布延迟消息. 不带上下文传递
func (ch *Channel) PublishDelayWithoutContext(ctx context.Context, exchange, key string, body []byte, timer time.Duration) (err error) {
	_, err = ch.Channel.PublishWithDeferredConfirmWithContext(ctx, exchange, key, false, false,
		amqp.Publishing{ContentType: "text/plain", Body: body, Expiration: fmt.Sprintf("%d", timer.Milliseconds())})
	return err
}

// QueueDeclare 创建队列.
func (ch *Channel) QueueDeclare(name string) (err error) {
	_, err = ch.Channel.QueueDeclare(name, true, false, false, false, nil)
	return
}

// QueueDeclareWithDelay 创建延迟队列.
func (ch *Channel) QueueDeclareWithDelay(name, exchange, key string) (err error) {
	_, err = ch.Channel.QueueDeclare(name, true, false, false, false, amqp.Table{
		"x-dead-letter-exchange":    exchange,
		"x-dead-letter-routing-key": key,
	})
	return
}

// QueueBind 绑定队列.
func (ch *Channel) QueueBind(name, key, exchange string) (err error) {
	return ch.Channel.QueueBind(name, key, exchange, false, nil)
}

// NewConsumer 实例化一个消费者, 会单独用一个channel.
func (ch *Channel) NewConsumer(ctx context.Context, queue string, handler func(context.Context, []byte) error) error {
	deliveries, err := ch.Consume(queue, "", false, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("consume err: %v, queue: %s", err, queue)
	}

	for msg := range deliveries {
		select {
		case <-ctx.Done():
			_ = msg.Reject(true)
			return fmt.Errorf("context cancel")
		default:
		}

		customCtx, body, ctxErr := GetContextFromMessage(msg.Body)
		if ctxErr != nil {
			_ = msg.Reject(true)
			continue
		}

		err = handler(customCtx, body)
		if err != nil {
			_ = msg.Reject(true)
			continue
		}
		_ = msg.Ack(false)
	}

	return nil
}

func (cf *ChannelFactory) MakeObject(ctx context.Context) (*pool.PooledObject, error) {
	//fmt.Printf("make object\n")
	//fmt.Printf("activeConn", activeConn)
	if activeConn.IsClosed() {
		mutex.Lock()
		var cErr error
		activeConn, cErr = Dial(cf.mqConnStr)
		if cErr != nil {
			mutex.Unlock()
			return nil, cErr
		}
		ChannelPool.Clear(ctx)
		mutex.Unlock()
	}
	ch, err := activeConn.Channel()
	if err != nil {
		return nil, err
	}
	return pool.NewPooledObject(
		ch,
	), nil
}

func (cf *ChannelFactory) DestroyObject(ctx context.Context, object *pool.PooledObject) error {
	ch := object.Object.(*Channel)
	//fmt.Printf("destroy object\n")
	if ch.IsClosed() {
		return nil
	}
	return ch.Close()
}

func (cf *ChannelFactory) ValidateObject(ctx context.Context, object *pool.PooledObject) bool {
	//fmt.Printf("validate Object\n")
	if activeConn.IsClosed() {
		return false
	}
	ch := object.Object.(*Channel)
	if ch.IsClosed() {
		return false
	}
	return true
}

func (cf *ChannelFactory) ActivateObject(ctx context.Context, object *pool.PooledObject) error {
	//fmt.Printf("activate object\n")
	ch := object.Object.(*Channel)
	if activeConn.IsClosed() || ch.IsClosed() {
		return amqp.Error{}
	}
	return nil
}

func (cf *ChannelFactory) PassivateObject(ctx context.Context, object *pool.PooledObject) error {
	// do passivate
	//fmt.Printf("passivate object\n")
	ch := object.Object.(*Channel)
	if ch.IsClosed() {
		return amqp.Error{}
	}
	return nil
}

// 从 CtxData 中创建上下文对象
func GetContextFromData(ctxData []CtxData) context.Context {
	ctx := context.Background()
	// 倒序排列
	for i := len(ctxData) - 1; i >= 0; i-- {
		cValue := ctxData[i]
		ctx = context.WithValue(ctx, cValue.Key, cValue.Value)
	}

	return ctx
}

// 将上下文信息添加到消息体中
func AddContextToMessage(ctx context.Context, data []byte) ([]byte, error) {
	// 将上下文对象转换为 map[string]any
	ctxData := GetKeyValues(ctx)

	dataWithCtx := DataWithCtx{
		Ctx:  ctxData,
		Data: data,
	}

	dataWithContext, err := json.Marshal(dataWithCtx)
	if err != nil {
		return nil, err
	}
	return dataWithContext, nil
}

// 从消息体中解析出上下文信息
func GetContextFromMessage(body []byte) (context.Context, []byte, error) {
	var dataWithCtx DataWithCtx
	err := json.Unmarshal(body, &dataWithCtx)
	if err != nil {
		ctx := context.Background()
		return ctx, body, nil
	}

	if len(dataWithCtx.Ctx) == 0 || len(dataWithCtx.Data) == 0 {
		ctx := context.Background()
		return ctx, body, nil
	}

	// 从 map 中创建上下文对象
	ctx := GetContextFromData(dataWithCtx.Ctx)

	return ctx, dataWithCtx.Data, nil
}
