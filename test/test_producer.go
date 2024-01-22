package main

import (
	"log"
	"time"

	"github.com/s290305915/go-rbmq/rbmq"
)

type OrderRbmqPorducer struct {
	*rbmq.RbmqInstance
}

func LoadProducer() *OrderRbmqPorducer {
	rabbit, err := rbmq.NewRabbit(&mqConf)
	if err != nil {
		return nil
	}
	// 生产者注册RabbitMQ
	orderProducerConfig := new(rbmq.ConsumerConfig)
	orderProducerConfig.ExchangeName = "test_exchange"
	orderProducerConfig.QueueName = "queue1"
	orderProducerConfig.KeyName = "key_consumer"
	orderProducerConfig.ExchangeType = rbmq.DIRECT_EXCHANGE
	//orderProducer := orderProducerConfig.NewInstance()
	orderProducer := orderProducerConfig.NewInstanceByConn(rabbit)

	return &OrderRbmqPorducer{
		RbmqInstance: orderProducer,
	}
}

func (c *OrderRbmqPorducer) Send(data []byte) {
	log.Println("start publisher:", c.ExchangeName, c.KeyName, string(data))
	go func() {
		err := c.MqChan.Publish(c.ExchangeName, c.KeyName, data)
		if err != nil {
			log.Fatalf("publish msg err: %v", err)
		}
	}()

	time.Sleep(1 * time.Second)
}
