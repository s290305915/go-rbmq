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

	cFun()

	err := orderProdc.Send(ctx, []byte(randomString(20)))
	if err != nil {
		t.Error(err)
	}

	fmt.Println("测试结束")

}
