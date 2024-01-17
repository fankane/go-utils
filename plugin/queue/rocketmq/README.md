### rocketmq
1. 在入口，比如 main.go 里面隐式导入kafka包路径
```go 
import _ "github.com/fankane/go-utils/plugin/queue/rocketmq"
```

2. 在运行文件根目录下的 **system_plugin.yaml** 文件(没有则新建一个)里面添加如下内容
```yaml
plugins:
  queue:  # 插件类型
    rocketmq:
      producers:
        default:                     # producer 名
          name_server:
            - 192.168.99.45:9876
          name_space: "hufan_ns"
          group_name: "test_group"
          retries: 3
          send_msg_timeout_ms: 10000 # 发送超时，单位：ms
        p2:
          name_server:
            - 192.168.99.45:9876
          name_space: "hufan_ns2"
          group_name: "test_group2"
          retries: 3
          send_msg_timeout_ms: 10000 # 发送超时，单位：ms
      consumers:
        c1:                          # consumer 名称
          topics:
            - test_topic_hf3
          async_consume: false #异步消费
          name_server:
            - 192.168.99.45:9876
          name_space: "hufan_ns"
          group_name: "tt"
          consume_from: 2
          consume_timestamp: "20240117144500"
          filter_history_for_init: true
```

3. 使用方式
- 3.1 生产者使用
```go
// 默认生产者发送消息
DefaultProducer.SendSync(context.Background(), "topicName", []byte("hello"))

// 获取配置文件里面 p2 对应的生产者发送消息
GetProducer("p2").SendSync(context.Background(), "topicName", []byte("hello"))
```
- 3.2 消费者使用
```go
// 给配置文件里面 c1 对应的消费者注册handler方法
// 当有消息过来时，会调用注册的 function并执行

RegisterHandler("c1", func(ctx context.Context, value []byte) error {
		fmt.Println("value:", string(value))
		return nil
	})
```

4. ConsumeFrom说明
- ConsumeFromLastOffset, ConsumeFromFirstOffset, ConsumeFromTimestamp
    > 只在 ConsumeGroup 第一次启动时有效, 后续再启动接着上次消费的进度开始消费
- 如果想第一次启动 consume Group 的时候，也能过滤历史信息，可以配置 
   ```filter_history_for_init: true```
    > 此配置，根据borntimestamp, 来过滤，<br >1. 如果是 ConsumeFromLastOffset，则过滤掉创建消费者之前的信息<br >2. 如果是 ConsumeFromTimestamp，则过滤掉consume_timestamp指定时间之前的信息，如果没有配置consume_timestamp，默认consume_timestamp为创建消费者前30分钟


 