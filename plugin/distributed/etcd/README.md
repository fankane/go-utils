### etcd
1. 在入口，比如 main.go 里面隐式导入etcd包路径
```go 
import _ "github.com/fankane/go-utils/plugin/distributed/etcd"
```

2. 在运行文件根目录下的 **system_plugin.yaml** 文件(没有则新建一个)里面添加如下内容
```yaml
plugins:
  distributed:  # 插件类型: 
    etcd: # 插件名
      default:                # etcd连接名称：default，可以是其他名字
        endpoints: ["192.168.0.93:2379"]
        username: xxx
        password: xxx
        dial_timeout_ms: 1000 # 连接超时时间，单位：毫秒，不填默认 5000ms
        open_discovery: true  # 开启服务发现功能
        server_info:          # open_discovery = true时必填
          server_name: "abc"  # [必填] 服务名
          server_id: "abc001" # [必填] 同一个服务，不同节点的ID，是服务的唯一标识
          region: wuhan       # [可选] 服务区域 
          host: 10.10.10.11   # [可选] 
          check_interval: 2   # [可选] 服务上报间隔时间，单位：秒，不填默认 3s

```

3. 在需要使用的地方，直接使用
```go
// 使用默认 default 的 etcd, 直接如下
etcd.Op.Get()

// 使用指定etcd连接, 如下:
etcd.GetOperate("xxx").Get()

// 监听服务信息, 每次服务有变动，serverChan 会收到最新的全量服务数据信息
// 特别注意：
//    只能获取到配置文件里面，server_info 里面配置的服务名对应的信息
//    比如服务名：abc有3个节点，ID分别为 abc001, abc002, abc003
//    1. 当有第4个服务abc004加入时，serverChan 会有全量的数据写入 [len(map) = 4]
//    2. 当abc001下线的时候，serverChan 会写入最新的全量数据[len(map) = 2]
//    3. 如果有其他的服务比如 xyz 的节点更新时，serverChan 不会有动静

serverChan := make(chan map[string]*ServerInfo)
go etcd.Op.WatchServers(serverChan)
for {
	serverInfoMap := <- serverChan
	doBusiness()
}

```