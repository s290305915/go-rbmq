package test

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/s290305915/go-rbmq/rbmq"
)

func BenchmarkProducer(b *testing.B) {
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

	orderProdc := LoadProducer()
	for i := 0; i < totalRequests; i++ {
		startTime := time.Now()
		wg.Add(1)

		//doMiter(i)

		go func(i int) {
			defer wg.Done()
			orderProdc.Send([]byte(randomString(20)))
		}(i)

		if (i+1)%maxConcurrency == 0 {
			counter++
			wg.Wait()
			endTime := time.Now()
			costTime := endTime.Sub(startTime).Seconds()
			fmt.Println("当前并发次数:", counter, "总请求数量:", i, "耗时:", costTime, "秒")
			//time.Sleep(1 * time.Second)
		}
	}

	wg.Wait()

	fmt.Println("测试结束")

}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func randomString(length int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, length)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
