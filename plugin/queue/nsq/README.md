### nsq
1. 在入口，比如 main.go 里面隐式导入kafka包路径
```go 
import _ "github.com/fankane/go-utils/plugin/queue/nsq"
```

2. 在运行文件根目录下的 **system_plugin.yaml** 文件(没有则新建一个)里面添加如下内容
```yaml
plugins:
  queue:  # 插件类型
    nsq:
      producers:
        default:                     # producer 名
          send_type: sync            # producer 发生类型 [sync, async]
          addr: 192.168.0.93:4150
        p2:
          send_type: async
          addr: 192.168.0.93:4150
      consumers:
        c1:                          # consumer 名称
          addrs:
            - 192.168.0.93:4161
          topic: test1
          channel: test_group       
          concurrency_consume: false   # 并发消费
          concurrency_max: 100         # 并发数，当concurrency_consume=true时有效，不填默认1000
        c2:                          # consumer 名称
          addrs:
            - 192.168.0.93:4161
          topic: test1
          channel: test_group2        
          concurrency_consume: false   # 并发消费
          concurrency_max: 100         # 并发数，当concurrency_consume=true时有效，不填默认1000
```

3. 使用方式
- 3.1 生产者使用
```go
// 默认生产者发送消息
DefaultProducer.SendMessage("topic name", []byte("value"))

// 获取配置文件里面 p2 对应的生产者发送消息
GetProducer("p2").SendMessage("topic name",  []byte("value"))
```
- 3.2 消费者使用
```go
// 给配置文件里面 c1 对应的消费者注册handler方法
// 当有消息过来时，会调用注册的 function并执行

RegisterHandler("c1", func(ctx context.Context, value []byte) error {
    fmt.Println(fmt.Sprintf("Attempts:%d, NSQDAddress:%s, time:%d", Attempts(ctx),
    NSQDAddress(ctx), Timestamp(ctx)), "value:", string(value))
    return nil
})
```