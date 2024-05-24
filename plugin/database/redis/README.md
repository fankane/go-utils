### Redis
1. 在入口，比如 main.go 里面隐式导入redis包路径
```go 
import _ "github.com/fankane/go-utils/plugin/database/redis"
```

2. 在运行文件根目录下的 **system_plugin.yaml** 文件(没有则新建一个)里面添加如下内容
```yaml
plugins:
  database:  # 插件类型: 
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
      client2:                # 连接名称：default，可以是其他名字
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
```

3. 在需要使用的地方，直接使用
    ```go
    使用默认 default 的log, 直接如下
    redis.Client.Get()
    
    使用指定redis连接, 如下:
    redis.GetClient("client2").Get()
    ```

4. Redis分布式锁
    > 基于lua脚本实现
    ```go
    lock := NewRdsLock(redis.Client)
    lock.Lock("key") // 加锁
    lock.Release() // 释放锁
    
    ```