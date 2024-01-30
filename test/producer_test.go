package test

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"

	"github.com/s290305915/go-rbmq/rbmq"
)

func TestProducerSingle(t *testing.T) {
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

	orderProdc := LoadProducer(mqConf)
	ctx := context.Background()
	ctx = context.WithValue(ctx, "ProductLine", "ggdxx")
	ctx = context.WithValue(ctx, "SaasToken", "--xxx--")
	ctx = context.WithValue(ctx, "AppId", "ffee_pc")

	ctx = context.WithValue(ctx, "ddd", "ggg")
	ctx = context.WithValue(ctx, "ProductLine_key", map[string]any{"key1": "ggg", "key2": 11199, "key3": struct {
		aa string
		bb int
	}{
		aa: "cc",
		bb: 180,
	}})

	cFun := func() { fmt.Println("测试结束") }
	ctx, cFun = context.WithTimeout(ctx, 1*time.Second)
	ctx = context.WithValue(ctx, "ddd", "123")

	go func() {
		for {
			// 发送消息带上下文
			err := orderProdc.Send(ctx, []byte("雷猴"))
			if err != nil {
				log.Fatalf("publish msg err: %v", err)
				return
			}

			// 发送消息不带上下文
			err = orderProdc.SendWithoutContext([]byte(randomString(50)))
			if err != nil {
				log.Fatalf("publish msg err(no context): %v", err)
				return
			}

			time.Sleep(1 * time.Second)
		}
	}()

	cFun()
	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)

	<-sigterm // Await a sigterm signal before safely closing the consumer
	log.Println("app shut down")

	fmt.Println("测试结束")

}
