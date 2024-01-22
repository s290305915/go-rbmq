package main

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/s290305915/go-rbmq/rbmq"
)

func BenchmarkConsumerTest(b *testing.B) {
	b.StopTimer()
	fmt.Println("启动测试进程")

	var wg sync.WaitGroup
	maxConcurrency := 10000 // 最大并发数量
	totalRequests := 100000 // 总请求数量

	//maxConcurrency := b.N // 最大并发数量
	//totalRequests := b.N  // 总请求数量

	counter := 0

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

	b.StartTimer()

	orderProdc := LoadConsumer(mqConf)
	for i := 0; i < totalRequests; i++ {
		startTime := time.Now()
		wg.Add(1)

		//doMiter(i)

		go func(i int) {
			defer wg.Done()
			orderProdc.Consume()
		}(i)

		if (i+1)%maxConcurrency == 0 {
			counter++
			wg.Wait()
			endTime := time.Now()
			costTime := endTime.Sub(startTime).Seconds()
			fmt.Println("当前并发次数:", counter, "总消费数量:", i, "耗时:", costTime, "秒")
			//time.Sleep(1 * time.Second)
		}
	}

	wg.Wait()

	fmt.Println("测试结束")

}
