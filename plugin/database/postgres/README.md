### Postgres
1. 在入口，比如 main.go 里面隐式导入postgres包路径
```go 
import _ "github.com/fankane/go-utils/plugin/database/postgres"
```

2. 在运行文件根目录下的 **system_plugin.yaml** 文件(没有则新建一个)里面添加如下内容
```yaml
plugins:
  database:  # 插件类型: 
    postgres: # 插件名
        default:                # 连接名称：default，可以是其他名字
          host: 127.0.0.1
          port: 5432
          user: root
          pwd: 1234
          db_name: test               # postgres 连接名称：default，可以是其他名字
          conn_max_life_time_sec: 600  # 连接最大存活时间, 不填连接建立后一直不会关闭 单位：秒
          conn_max_idle_time_sec: 600  # 空闲连接最大存活时间, 单位：秒
          max_open_conn: 10            # 最大连接数, 不填默认无限制
          max_idle_conn: 2             # 最大空闲连接数，不填，mysql1.7.1 默认为2
        name2:                # 连接名称：default，可以是其他名字
          host: 192.168.0.1
          port: 5432
          user: root
          pwd: 1234
          db_name: test               # postgres 连接名称：default，可以是其他名字
          conn_max_life_time_sec: 600  # 连接最大存活时间, 不填连接建立后一直不会关闭 单位：秒
          conn_max_idle_time_sec: 600  # 空闲连接最大存活时间, 单位：秒
          max_open_conn: 10            # 最大连接数, 不填默认无限制
          max_idle_conn: 2             # 最大空闲连接数，不填，mysql1.7.1 默认为2

```

3. 在需要使用的地方，直接使用
```go
使用默认 default 的log, 直接如下
postgres.DB.Query()

使用指定postgres连接, 如下:
postgres.GetDB("name2").Query()
```