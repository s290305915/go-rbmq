# Welcome To go-rbmq

<a href="https://github.com/s290305915">
    <img src="https://badgen.net/badge/Github/s290305915?icon=github" alt="">
</a>
<img alt="GitHub last commit" src="https://img.shields.io/github/last-commit/s290305915/go-rmbq">
<img alt="GitHub Repo stars" src="https://img.shields.io/github/stars/s290305915/go-rmbq?style=social">
</p>

## 简介

rabbit-mq 快速接入程序，本项目参考自[shixiaofeia/fly](https://github.com/shixiaofeia/fly)项目所作二次开发，主要目的在于快速接入rabbit-mq推送和消费程序，引入即可使用，无需花时间在基础组件

## 快速接入

```
go get github.com/s290305915/go-rbmq
```


## 项目结构

```
|── rbmq                // rabbit-mq主要操作
|  ├── base             // 初始化
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

[amqp091-go](github.com/rabbitmq/amqp091-go)