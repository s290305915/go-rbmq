package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/s290305915/go-rbmq/rbmq"
)

func Test() {

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

	orderConsu := LoadConsumer(mqConf)
	go orderConsu.Consume()

	//orderProdc := LoadProducer()
	//
	//go func() {
	//	for {
	//		orderProdc.Send([]byte("雷猴"))
	//	}
	//}()

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)

	<-sigterm // Await a sigterm signal before safely closing the consumer
	log.Println("app shut down")
}
