# go-utils
实际项目开发时，常用的工具方法汇总。

> 如果你需要什么常见的功能，这里没有的，可提起pull request, 或者联系本人添加 <br>
> 邮箱: fanhu1116@qq.com 

## plugin【插件使用】
#### <font style="color: red">不需要写代码</font> 将一些基础功能，通过配置文件的形式，在服务启动的时候，自动加载，需要的时候，直接使用； <br>

```go
import	(
    "github.com/fankane/go-utils/plugin"
	
    // 按需引入即可
    _ "github.com/fankane/go-utils/plugin/log"
    _ "github.com/fankane/go-utils/plugin/queue/memory"
)

func main() {
    plugin.Load()
}

```

#### [较完整的插件配置文件](./plugin/README.md) 

详细说明文档
- [log 插件使用](./plugin/log/README.md)
- 监控
  - [Prometheus 插件使用](./plugin/monitor/prometheus/README.md)
- Database
  - [MySQL 插件使用](./plugin/database/mysql/README.md)
  - [Postgres 插件使用](./plugin/database/postgres/README.md)
  - [Oracle 插件使用](./plugin/database/oracle/README.md)
  - [Freecache插件使用](./plugin/database/freecache/README.md)
  - [Redis 插件使用](./plugin/database/redis/README.md)
  - [Elasticsearch 插件使用](./plugin/database/es/README.md)
  - [Neo4j 插件使用](./plugin/database/neo4j/README.md)
  - [InfluxDB 插件使用](./plugin/database/influx/README.md)
- 服务
  - [Pprof 插件使用](./plugin/serve/pprof/README.md)
  - [Conf 插件使用](./plugin/serve/conf/README.md)
  - [Nacos 插件使用](./plugin/serve/nacos/README.md)
- 队列
  - [Kafka 插件使用](./plugin/queue/kafka/README.md)
  - [NSQ 插件使用](plugin/queue/nsq/README.md)
  - [自研内存队列(支持延时)](plugin/queue/memory/README.md)
  - [Rabbit 插件使用](plugin/queue/rabbit/README.md)
  - [Rocketmq 插件使用](plugin/queue/rocketmq/README.md)
- 分布式
  - [etcd 插件使用](plugin/distributed/etcd/README.md)
  - [jaeger 插件使用](plugin/distributed/jaeger/README.md)


## 常用工具
### 函数执行
  - 多次运行 【可选：重试次数，超时时间，间隔时间】
  - 单次运行 【加锁，执行耗时，最长等待】
  - 并发执行 【协程池，任务管理器】
  - 定时执行、延迟执行、ticker执行
### 文件操作
  - csv, xls, xlsx 读写
  - rar, zip 读写
  - 文件、目录 读写 操作
### 数据结构
  - 类型断言 【数字型，数组/切片型】
  - 类型转换 【slice, string, float, int, interface, bytes 等之间的转换】
  - 切片元素包含关系
  - JSON 字符串转 类JSONSchema
  - 中文字符串处理【长度获取，截取子串】
  - error 封装：支持code、msg、showMsg ;同时也实现了Error方法，可跟原生 error 兼容使用
### 硬件信息
  - CPU, 内存信息, 磁盘信息 【机器整体、指定进程、当前进程】
### 网络
  - Http 封装
  - [Websocket封装](http/README.md)
  - SSE[server sent events] 封装
