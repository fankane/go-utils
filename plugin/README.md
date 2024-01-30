## 较完整的插件用例

使用插件
- 使用 log
- 使用2个MySQL
- 使用1个Redis
- 使用1个postgres
- 使用pprof
- 使用conf
- 使用 Prometheus

```yaml
plugins:
  log:  # 插件类型: log
    zap_sugar: # 插件名
      default:                # 日志名称：default，可以是其他名字
        level: debug          # 日志级别, 默认info [debug, info, warn, error, panic]
        enable_stdout: true   # 是否开启日志同步到控制台, 默认false
        filename: ./test.log  # 日志文件，不填默认使用 ./log.log
        max_size: 1           # 日志文件滚动日志的大小 单位 MB
        max_age: 2            # 最大日志保留天数
        max_backups: 7        # 最大日志备份数量
        compress: false       # 是否压缩，默认：false 

  database:  # 插件类型
    mysql: # 插件名
      default:                # MySQL连接名称：default，可以是其他名字
        dsn: user:pwd@tcp(127.0.0.1:3306)/dbName?parseTime=true
        conn_max_life_time_sec: 600  # 连接最大存活时间, 不填连接建立后一直不会关闭 单位：秒
        conn_max_idle_time_sec: 600  # 空闲连接最大存活时间, 单位：秒
        max_open_conn: 10            # 最大连接数, 不填默认无限制
        max_idle_conn: 2             # 最大空闲连接数，不填，mysql1.7.1 默认为2
      mysql2:                        # 
        dsn: user:pwd@tcp(127.0.0.1:3306)/dbName?parseTime=true
        conn_max_life_time_sec: 600  # 连接最大存活时间, 不填连接建立后一直不会关闭 单位：秒
        conn_max_idle_time_sec: 600  # 空闲连接最大存活时间, 单位：秒
        max_open_conn: 10            # 最大连接数, 不填默认无限制
        max_idle_conn: 2             # 最大空闲连接数，不填，mysql1.7.1 默认为2
    redis: # 插件名
      default:                # 连接名称：default，可以是其他名字
        addr: 127.0.0.1:6379
        db: 0                        # 不填默认为0
        user: xx
        pwd: xx
        dial_timeout_mils: 0         #建立连接超时时间，单位：毫秒，
        ping_timeout_mils: 0         #创建client后ping测试超时时间，单位：毫秒，不填默认为1000ms
        conn_max_life_time_sec: 600  # 连接最大存活时间, 不填连接建立后一直不会关闭 单位：秒
        conn_max_idle_time_sec: 600  # 空闲连接最大存活时间, 单位：秒
        min_idle_conn: 10            # 最小空闲连接数,
        max_idle_conn: 2             # 最大空闲连接数
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

  serve:  # 插件类型: 服务类
    pprof: # 插件名
      addr: 127.0.0.1:6060
    conf: # 插件名
      conf_file: test_conf.yaml
      watch_change: true           #监听文件更新, 默认false
      change_cron: "@every 10s"

  monitor:  # 插件类型:
    prometheus: # 插件名
      port: 7701
      path: "/metrics"
      custom_collects:
        - coll_type: counter     # 采集类型[counter, gauge, histogram, summary]
          info:
            counter_test1:       # 指标名
              help: 自定义计数指标1 # 指标说明
              labels:
                - label1    # 标签
                - label2    # 标签
            counter_test2:
              help: 自定义计数指标2
              labels:
                - label1
                - label2
        - coll_type: gauge  
          info:
            test1:
              help: 自定义数值指标1
              labels:
                - label1
                - label2
            test2:
              help: 自定义数值指标1
              labels:
                - label1
                - label2
```
