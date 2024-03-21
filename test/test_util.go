package test

import (
	"context"
	"log"
	"math/rand"
	"time"

	"github.com/s290305915/go-rbmq/rbmq"
)

type OrderRbmqConsumer struct {
	*rbmq.RbmqInstance
}

func LoadConsumer2(mqConf rbmq.Conf) *OrderRbmqConsumer {
	rabbit, err := rbmq.NewRabbit(&mqConf)
	if err != nil {
		return nil
	}

	// 消费者注册RabbitMQ
	orderConsumerConfig := new(rbmq.ConsumerConfig)
	orderConsumerConfig.ExchangeName = "test_exchange"
	orderConsumerConfig.QueueName = "queue1"
	orderConsumerConfig.KeyName = "key_consumer"
	orderConsumerConfig.ExchangeType = rbmq.DIRECT_EXCHANGE

	//orderConsumer := orderConsumerConfig.NewInstance()
	//log.Printf("orderConsu is %+v", orderConsumer)
	consumer := orderConsumerConfig.NewInstanceByConn(rabbit)
	return &OrderRbmqConsumer{
		RbmqInstance: consumer,
	}
}

func (c *OrderRbmqConsumer) Consume() {
	//fmt.Printf("c is %+v", c)
	conf := c.ConsumerConfig
	log.Printf("start consumer ok: %s, %s, %s\n", conf.ExchangeName, conf.QueueName, conf.KeyName)

	go func() {
		err := c.MqChan.NewConsumer(context.Background(), conf.QueueName, func(ctx context.Context, body []byte) error {
			log.Println("consumer messages ------------------> :", string(body))

			log.Printf("consumer messages with context-value:%+v \n", ctx)

			prd1 := ctx.Value("ProductLine")
			prd2 := ctx.Value("SaasToken")
			prd3 := ctx.Value("AppId")
			prd4 := ctx.Value("ProductLine_key")

			log.Printf("ProductLine:%+v \n", prd1)
			log.Printf("SaasToken:%+v \n", prd2)
			log.Printf("AppId:%+v \n", prd3)
			log.Printf("ProductLine_key:%+v \n", prd4)

			return nil
		})
		if err != nil {
			log.Fatalf("consume err: %v", err)
		}
	}()
}

type OrderRbmqPorducer struct {
	*rbmq.RbmqInstance
}

func LoadProducer(mqConf rbmq.Conf) *OrderRbmqPorducer {
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

func (c *OrderRbmqPorducer) Send(ctx context.Context, data []byte) error {
	log.Println("start publisher:", c.ExchangeName, c.KeyName, string(data))
	go func() {
		pErr := c.MqChan.Publish(ctx, c.ExchangeName, c.KeyName, data)
		if pErr != nil {
			log.Fatalf("publish msg err: %v", pErr)
		}
	}()

	time.Sleep(1 * time.Second)
	return nil
}

func (c *OrderRbmqPorducer) SendTx(ctx context.Context, data []byte) (*rbmq.Channel, error) {
	log.Println("start publisher with transaction:", c.ExchangeName, c.KeyName, string(data))
	ch, pErr := c.MqChan.PublishTx(ctx, c.ExchangeName, c.KeyName, data)
	if pErr != nil {
		log.Fatalf("publish msg err: %v", pErr)
	}

	time.Sleep(1 * time.Second)
	return ch, nil
}

func (c *OrderRbmqPorducer) SendWithoutContext(data []byte) error {
	log.Println("start publisher(no context):", c.ExchangeName, c.KeyName, string(data))
	go func() {
		pErr := c.MqChan.PublishWithoutContext(context.Background(), c.ExchangeName, c.KeyName, data)
		if pErr != nil {
			log.Fatalf("publish msg err(no context): %v", pErr)
		}
	}()

	time.Sleep(1 * time.Second)
	return nil
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func randomString(length int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, length)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
