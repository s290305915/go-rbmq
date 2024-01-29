# Welcome To go-rbmq

<a href="https://github.com/s290305915">
    <img src="https://badgen.net/badge/Github/s290305915?icon=github" alt="">
</a>
<img alt="GitHub last commit" src="https://img.shields.io/github/last-commit/s290305915/go-rmbq">
<img alt="GitHub Repo stars" src="https://img.shields.io/github/stars/s290305915/go-rmbq?style=social">

## 简介

rabbit-mq 快速接入程序，本项目参考自[shixiaofeia/fly](https://github.com/shixiaofeia/fly)项目所作二次开发，主要目的在于快速接入rabbit-mq推送和消费程序，引入即可使用，无需花时间在基础组件

## 快速使用

```
go get github.com/s290305915/go-rbmq
```

设置配置文件：
```json
{
    "addr": "127.0.0.1",
    "port": "5672",
    "user": "admin",
    "pwd": "123",
    "vhost": "/",
    "pool_idle":{
        "max_size": 200,
        "min_idle": 50,
        "max_idle": 2000
    }
}
```


## 项目结构

```
|── rbmq                // rabbit-mq主要操作
|  ├── base             // 初始化
|  ├── context_mapping  // 上下文获取
|  ├── exchange_types   // 交换机类型枚举
|  ├── instance_config  // 接入各类参数
|  ├── instance         // 实例化
|  └── rabbitmq         // 基础设施
├── test                // 消费者示例
|  ├── test_consumer    // 消费者实例代码（可直接复制使用，根据业务修改参数即可）
|  ├── test_producer    // 消费者实例代码（可直接复制使用，根据业务修改参数即可）
├── go.mod              // 包管理    
├── go.sum              // 包管理    
├── README.md
```

## 项目引用

### amqp091-go
### go-commons-pool 用于实现连接池处理

## 2023年12月1日 更新
增加go-commons-pool连接池配置

## 2021年1月5日 更新
增加context_mapping上下文获取并传递到消息和解析还原上下文

生产者使用方法详见：[test_consumer](./test/producer_test.go)

消费者使用方法详见：[test_consumer](./test/consumer_test.go)

消费者需添加参数签名`func(ctx context.Context, body []byte) error`，其中ctx为上下文，body为消息体

[amqp091-go](github.com/rabbitmq/amqp091-go)