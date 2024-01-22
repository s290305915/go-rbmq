package test

import (
	"context"
	"log"

	"github.com/s290305915/go-rbmq/rbmq"
)

type OrderRbmqConsumer struct {
	*rbmq.RbmqInstance
}

func LoadConsumer() *OrderRbmqConsumer {

	// 消费者注册RabbitMQ
	orderConsumerConfig := new(rbmq.ConsumerConfig)
	orderConsumerConfig.ExchangeName = "test_exchange"
	orderConsumerConfig.QueueName = "queue1"
	orderConsumerConfig.KeyName = "key_consumer"
	orderConsumerConfig.ExchangeType = rbmq.DIRECT_EXCHANGE

	orderConsumer := orderConsumerConfig.NewInstance()
	log.Printf("orderConsu is %+v", orderConsumer)

	return &OrderRbmqConsumer{
		RbmqInstance: orderConsumer,
	}
}

func (c *OrderRbmqConsumer) Consume() {
	//fmt.Printf("c is %+v", c)
	conf := c.ConsumerConfig
	log.Printf("start consumer: %s, %s, %s\n", conf.ExchangeName, conf.QueueName, conf.KeyName)

	go func() {
		err := c.MqChan.NewConsumer(context.Background(), conf.QueueName, func(body []byte) error {
			log.Println("consumer messages ------------------> :", string(body))

			return nil
		})
		if err != nil {
			log.Fatalf("consume err: %v", err)
		}
	}()
}
