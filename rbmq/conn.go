package rbmq

import (
	"fmt"
	"sync"

	"github.com/rabbitmq/amqp091-go"
)

type Rabbit struct {
	config  *Conf
	conn    *amqp091.Connection
	channel *amqp091.Channel
	sync.Mutex
}

// NewServiceRabbit create a Rabbit by service rabbit config
// Please make sure your service has rabbit config
func NewServiceRabbit(serviceName string) (*Rabbit, error) {
	// TODO 通过服务名获取配置
	//config, err := GetServiceConnConfig(serviceName)
	//if err != nil {
	//	return nil, err
	//}

	//return NewRabbit(config)
	return nil, nil
}

// NewRabbit create Rabbit from ConnConfig
func NewRabbit(config *Conf) (*Rabbit, error) {
	// TODO 后续可在此处接入监控指标
	// OptionMetrics()
	conn, err := newConn(config)
	if err != nil {
		return nil, err
	}

	return &Rabbit{
		config: config,
		conn:   conn,
	}, err
}

const connStrFormat = "amqp://%v:%v@%v:%v/%v"

func newConn(config *Conf) (*amqp091.Connection, error) {
	vhost := config.Vhost
	if vhost == "" {
		vhost = "/"
	}
	connStr := fmt.Sprintf(connStrFormat, config.User, config.Pwd, config.Addr, config.Port, vhost)
	return amqp091.Dial(connStr)
}
