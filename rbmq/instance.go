package rbmq

import (
	"context"
	"log"
)

type RbmqInstance struct {
	*ConsumerConfig
	MqChan *Channel
}

func (c *ConsumerConfig) NewInstance() *RbmqInstance {
	// 方法实现

	//fmt.Printf("ChannelPool：%+v", ChannelPool)

	obj, err := ChannelPool.BorrowObject(context.TODO())
	if err != nil {
		return nil
	}
	ch := obj.(*Channel)

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

func (c *ConsumerConfig) NewInstanceByConn(r *Rabbit) *RbmqInstance {
	defaultConn := &Connection{r.conn}
	channel, err := defaultConn.Channel()
	if err != nil {
		log.Fatalf("create channel err: %v", err)
	}

	if err := channel.ExchangeDeclare(c.ExchangeName, c.ExchangeType.String()); err != nil {
		log.Fatalf("create exchange err: %v", err)
	}

	if err := channel.QueueDeclare(c.QueueName); err != nil {
		log.Fatalf("create queue err: %v", err)
	}

	if err := channel.QueueBind(c.QueueName, c.KeyName, c.ExchangeName); err != nil {
		log.Fatalf("bind queue err: %v", err)
	}

	log.Printf("init mq success, exchange: %s, queue: %s, key: %s", c.ExchangeName, c.QueueName, c.KeyName)
	return &RbmqInstance{
		ConsumerConfig: c,
		MqChan:         channel,
	}
}
