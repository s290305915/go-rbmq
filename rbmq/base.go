package rbmq

import (
	"context"
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type (
	Conf struct {
		Addr  string `json:"addr" yaml:"addr"`
		Port  string `json:"port" yaml:"port"`
		User  string `json:"user" yaml:"user"`
		Pwd   string `json:"pwd" yaml:"pwd"`
		Vhost string `json:"vhost" yaml:"vhost"`
	}
)

var (
	defaultConn    *Connection
	defaultChannel *Channel
)

// Init 初始化
func Init(c Conf) (ch *Channel, err error) {
	if c.Addr == "" {
		return nil, nil
	}
	defaultConn, err = Dial(fmt.Sprintf("amqp://%s:%s@%s:%s/%s",
		c.User,
		c.Pwd,
		c.Addr,
		c.Port,
		c.Vhost))
	if err != nil {
		return nil, fmt.Errorf("new mq conn err: %v", err)
	}

	defaultChannel, err = defaultConn.Channel()
	if err != nil {
		return nil, fmt.Errorf("new mq channel err: %v", err)
	}
	return defaultChannel, nil
}

// ExchangeDeclare 创建交换机.
func (ch *Channel) ExchangeDeclare(name string, kind string) (err error) {
	return ch.Channel.ExchangeDeclare(name, kind, true, false, false, false, nil)
}

// Publish 发布消息.
func (ch *Channel) Publish(exchange, key string, body []byte) (err error) {
	_, err = ch.Channel.PublishWithDeferredConfirmWithContext(context.Background(), exchange, key, false, false,
		amqp.Publishing{ContentType: "text/plain", Body: body})
	return err
}

// PublishWithDelay 发布延迟消息.
func (ch *Channel) PublishWithDelay(exchange, key string, body []byte, timer time.Duration) (err error) {
	_, err = ch.Channel.PublishWithDeferredConfirmWithContext(context.Background(), exchange, key, false, false,
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
func (ch *Channel) NewConsumer(ctx context.Context, queue string, handler func([]byte) error) error {
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
		err = handler(msg.Body)
		if err != nil {
			_ = msg.Reject(true)
			continue
		}
		_ = msg.Ack(false)
	}

	return nil
}
