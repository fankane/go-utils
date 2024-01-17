### kafka
1. 在入口，比如 main.go 里面隐式导入kafka包路径
```go 
import _ "github.com/fankane/go-utils/plugin/queue/kafka"
```

2. 在运行文件根目录下的 **system_plugin.yaml** 文件(没有则新建一个)里面添加如下内容
```yaml
plugins:
  queue:  # 插件类型
    kafka:
      producers:
        default:                     # producer 名
          send_type: sync            # producer 发生类型 [sync, async]
          addrs:
            - 192.168.0.93:9092
        p2:
          send_type: async
          addrs:
            - 192.168.0.93:9092
      consumers:
        c1:                          # consumer 名称
          addrs:
            - 192.168.0.93:9092
          topics:
            - test1
            - test2
          group_id: test_group       # 不填会默认生成一个随机值
          concurrency_consume: false   # 并发消费
          concurrency_max: 100         # 并发数，当concurrency_consume=true时有效，不填默认1000
          offset_initial: -1           # [可选] offset 初始值，当 consumer 时生效
          reset_offset_info:           # [可选] 设置使用 consumerGroup 消费时的初始offset值
            - topic: test1
              offset: 0              # -1 newest,
              set_for_all: true       # 是否对所有 partition 生效
            - topic: test2
              offset: -1
              set_for_all: false
              partitions_setting:     # 当 set_for_all = false 时生效
                - partition: 0
                  offset: 2
```

3. 使用方式
- 3.1 生产者使用
```go
// 默认生产者发送消息
DefaultProducer.SendMessage("topic name", []byte("key")), []byte("value"))

// 获取配置文件里面 p2 对应的生产者发送消息
GetProducer("p2").SendMessage("topic name", []byte("key")), []byte("value"))
```
- 3.2 消费者使用
```go
// 给配置文件里面 c1 对应的消费者注册handler方法
// 当有消息过来时，会调用注册的 function并执行

RegisterHandler("c1", func(ctx context.Context, key, value []byte) error {
		fmt.Println(fmt.Sprintf("topic:%s, partition:%d, offset:%d, timestamp:%s", Topic(ctx), Partition(ctx), Offset(ctx), Timestamp(ctx)),
			"business key:", string(key), "value:", string(value))
		return nil
	})
```