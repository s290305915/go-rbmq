package main

import (
	"github.com/s290305915/go-rbmq/rbmq"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	mqConf := rbmq.Conf{
		Addr:  "127.0.0.1",
		Port:  "5672",
		User:  "admin",
		Pwd:   "123",
		Vhost: "/",
	}

	time.Sleep(3 * time.Second)

	orderConsu := LoadConsumer(mqConf)
	orderConsu.Consume()

	//orderProdc := LoadProducer(mqConf)

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
