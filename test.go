package test

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/s290305915/go-rbmq/rbmq"
	"github.com/s290305915/go-rbmq/test_consumer"
	"github.com/s290305915/go-rbmq/test_producer"
)

func test() {

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
