package test

import (
	"context"
	"fmt"
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
	ctx = context.WithValue(ctx, "ProductLine", "PD1")
	ctx = context.WithValue(ctx, "SaasToken", "")
	ctx = context.WithValue(ctx, "AppId", "PD1_APP")

	ctx = context.WithValue(ctx, "ddd", "123")
	ctx = context.WithValue(ctx, "ProductLine_key", map[string]any{"key1": "123", "key2": 99, "key3": struct {
		aa string
		bb int
	}{
		aa: "123",
		bb: 99,
	}})

	cFun := func() { fmt.Println("测试结束") }
	ctx, cFun = context.WithTimeout(ctx, 1*time.Second)
	ctx = context.WithValue(ctx, "ddd", "123")

	cFun()

	err := orderProdc.Send(ctx, []byte(randomString(20)))
	if err != nil {
		t.Error(err)
	}

	fmt.Println("测试结束")

}
