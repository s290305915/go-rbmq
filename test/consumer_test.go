package test

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"

	"github.com/s290305915/go-rbmq/rbmq"
)

func TestConsumerSingle(t *testing.T) {
	fmt.Println("启动测试进程")

	mqConf := rbmq.Conf{
		Name:  "rbmq1",
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
	_ = rbmq.Conf{
		Name:  "rbmq2",
		Addr:  "127.0.0.2",
		Port:  "5672",
		User:  "admin",
		Pwd:   "123",
		Vhost: "cxl_host",
		PoolIdle: rbmq.Idle{
			MaxSize: 100,
			MinIdle: 1000,
			MaxIdle: 2000,
		},
	}

	time.Sleep(3 * time.Second)
	rbmq.Init(mqConf)

	orderProdc := LoadConsumer2(mqConf)
	go orderProdc.Consume()

	//rbmq.Init(mqConf2)
	//orderProdc2 := LoadConsumer2(mqConf2)
	//go orderProdc2.Consume()

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)

	<-sigterm // Await a sigterm signal before safely closing the consumer
	log.Println("app shut down")

	fmt.Println("测试结束")
}
