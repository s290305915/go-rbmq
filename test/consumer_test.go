package test

import (
	"fmt"
	"testing"
	"time"

	"github.com/s290305915/go-rbmq/rbmq"
)

func TestConsumerSingle(t *testing.T) {
	fmt.Println("启动测试进程")

	mqConf := rbmq.Conf{
		Addr:  "127.0.0.1",
		Port:  "5672",
		User:  "admin",
		Pwd:   "123",
		Vhost: "/",
		PoolIdle: rbmq.Idle{
			MaxSize: 100,
			MinIdle: 1000,
			MaxIdle: 2000,
		},
	}

	time.Sleep(3 * time.Second)
	rbmq.Init(mqConf)

	orderProdc := LoadConsumer2(mqConf)
	orderProdc.Consume()

	fmt.Println("测试结束")
}
