### Oracle
1. 在入口，比如 main.go 里面隐式导入 Oracle 包路径
```go 
import _ "github.com/fankane/go-utils/plugin/database/oracle"
```

2. 在运行文件根目录下的 **system_plugin.yaml** 文件(没有则新建一个)里面添加如下内容
```yaml
plugins:
  database:  # 插件类型
    oracle: # 插件名
      default:                # Oracle 连接名称：default，可以是其他名字
        host: 192.168.0.93
        port: 3306
        user: root
        pwd: 123456
        sid: helowin
        server_name:                 # 优先使用sid，没有则使用服务名
        conn_max_life_time_sec: 600  # 连接最大存活时间, 不填连接建立后一直不会关闭 单位：秒
        conn_max_idle_time_sec: 600  # 空闲连接最大存活时间, 单位：秒
        max_open_conn: 500            # 最大连接数, 不填默认无限制
        max_idle_conn: 2             # 最大空闲连接数，
      testName:                # Oracle 连接名称：default，可以是其他名字
        host: 192.168.0.93
        port: 3306
        user: root
        pwd: 123456
        sid: helowin
        server_name:                 # 优先使用sid，没有则使用服务名
        conn_max_life_time_sec: 600  # 连接最大存活时间, 不填连接建立后一直不会关闭 单位：秒
        conn_max_idle_time_sec: 600  # 空闲连接最大存活时间, 单位：秒
        max_open_conn: 500            # 最大连接数, 不填默认无限制
        max_idle_conn: 2             # 最大空闲连接数，
```

3. 在需要使用的地方，直接使用
```go
使用默认 default 的db, 直接如下
oracle.DB.Query()

使用指定MySQL连接, 如下:
oracle.GetDB("Oracle").Query()
```