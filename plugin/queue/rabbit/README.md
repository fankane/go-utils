### rabbit
1. 在入口，比如 main.go 里面隐式导入 rabbit 包路径
```go 
import _ "github.com/fankane/go-utils/plugin/queue/rabbit"
```

2. 在运行文件根目录下的 **system_plugin.yaml** 文件(没有则新建一个)里面添加如下内容
```yaml
plugins:
  queue:  # 插件类型
    rabbit:
      producers:
        default:                     # producer 名
          url: "amqp://guest:guest@192.168.0.93:5672/"
          durable: false
          auto_delete: false
          exclusive: false
          no_wait: false
        p2:
          url: "amqp://guest:guest@localhost:5672/"
      consumers:
        c1:                          # consumer 名称
          url: "amqp://guest:guest@192.168.0.93:5672/"
          durable: false
          auto_delete: false
          exclusive: false
          no_wait: false
          queue_names:
            - "hello -1"
        c2:                          # consumer 名称
          url: "amqp://guest:guest@192.168.0.93:5672/"

```

3. 使用方式
- 3.1 生产者使用
```go
// 默认生产者发送消息
DefaultProducer.SendMessage(ctx, "queue name", []byte("value"))

// 获取配置文件里面 p2 对应的生产者发送消息
GetProducer("p2").SendMessage(ctx, "queue name",  []byte("value"))
```
- 3.2 消费者使用
```go
// 给配置文件里面 c1 对应的消费者注册handler方法
// 当有消息过来时，会调用注册的 function并执行

RegisterHandler("c1", func(ctx context.Context, value []byte) error {
    fmt.Println(string(value))
    return nil
})
```