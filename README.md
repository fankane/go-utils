# go-utils
实际项目开发时，常用的工具方法汇总。

> 如果你需要什么常见的功能，这里没有的，可提起pull request, 或者联系本人添加 <br>
> 邮箱: fanhu1116@qq.com 

## plugin
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
  - [Freecache插件使用](./plugin/database/freecache/README.md)
  - [Redis 插件使用](./plugin/database/redis/README.md)
  - [Elasticsearch 插件使用](./plugin/database/es/README.md)
  - [Neo4j 插件使用](./plugin/database/neo4j/README.md)
- 服务
  - [Pprof 插件使用](./plugin/serve/pprof/README.md)
  - [Conf 插件使用](./plugin/serve/conf/README.md)
  - [Nacos 插件使用](./plugin/serve/nacos/README.md)
- 队列
  - [Kafka 插件使用](./plugin/queue/kafka/README.md)
  - [NSQ 插件使用](plugin/queue/nsq/README.md)
  - [自研内存队列](plugin/queue/memory/README.md)
  - [Rabbit 插件使用](plugin/queue/rabbit/README.md)
- 分布式
  - [etcd 插件使用](plugin/distributed/etcd/README.md)

## file
> 文件相关操作
- 文件、目录 读写 操作
- xlsx: xlsx 文件处理
- xls: xls 文件的读[建议优先考虑xlsx]
- csv: csv 文件处理

## archive
> 压缩文件相关操作
- rar: rar文件的相关操作
  - UnRar: 解压 rar 文件
- zip: zip文件的相关操作
  - CreateZip: 创建 zip 文件
  - UnZip: 解压 zip 文件

## string
> 字符串操作
  - 字符串转换、中文字符串处理
  - uuid
## utime
> 日期相关功能
  - cron 语法定期执行
  - 延期执行，ticker 执行，可选参数，比如最多执行次数，最长等待时间

## slice
- contain
  - InInterfaces,InStrings, InInts 等，判定某个切片里面是否存在某个具体的值; 支持基础类型 int, string, float
- transform
  - ToInterfaceSli: 将普通切片转换成interface切片，例如 []int -> []interface; 支持基础类型 int, string, float
- compare
  - StrSliContentEqual: 比较两个字符串切片内容是否相同，忽略顺序

## random

## uerr
> 自定义error，支持code、msg、showMsg ;同时也实现了Error方法，可跟原生 error 兼容使用

## http
> http 方法包装
- [websocket使用说明](http/README.md)