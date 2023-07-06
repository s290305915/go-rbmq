package rbmq

import (
	"log"
)

type RbmqInstance struct {
	*ConsumerConfig
	MqChan *Channel
}

func (c *ConsumerConfig) NewInstance(conf Conf) *RbmqInstance {
	// 方法实现

	ch, err := Init(conf)
	if err != nil {
		log.Fatal("init rabbit mq err: " + err.Error())
	}

	if err := ch.ExchangeDeclare(c.ExchangeName, c.ExchangeType.String()); err != nil {
		log.Fatalf("create exchange err: %v", err)
	}

	if err := ch.QueueDeclare(c.QueueName); err != nil {
		log.Fatalf("create queue err: %v", err)
	}

	if err := ch.QueueBind(c.QueueName, c.KeyName, c.ExchangeName); err != nil {
		log.Fatalf("bind queue err: %v", err)
	}

	log.Printf("init mq success, exchange: %s, queue: %s, key: %s", c.ExchangeName, c.QueueName, c.KeyName)
	return &RbmqInstance{
		ConsumerConfig: c,
		MqChan:         ch,
	}
}
