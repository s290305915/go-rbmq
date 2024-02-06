package rbmq

type ConsumerConfig struct {
	ExchangeName string
	QueueName    string
	KeyName      string
	ExchangeType EnumExchagne
	MqConf       Conf
}
