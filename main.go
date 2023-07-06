package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-rbmq/rbmq"
	"go-rbmq/test_consumer"
	"go-rbmq/test_producer"
)

func main() {

	mqConf := rbmq.Conf{
		Addr:  "-",
		Port:  "5672",
		User:  "-",
		Pwd:   "-",
		Vhost: "-",
	}

	time.Sleep(3 * time.Second)

	orderConsu := test_consumer.Load(mqConf)
	orderConsu.Consume()

	orderProdc := test_producer.Load(mqConf)

	go func() {
		for {
			orderProdc.Send([]byte("雷猴"))
		}
	}()

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)

	<-sigterm // Await a sigterm signal before safely closing the consumer
	log.Println("app shut down")
}
