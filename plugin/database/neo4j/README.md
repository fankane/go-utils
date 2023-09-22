### neo4j
1. 在入口，比如 main.go 里面隐式导入neo4j包路径
```go 
import _ "github.com/fankane/go-utils/plugin/database/neo4j"
```

2. 在运行文件根目录下的 **system_plugin.yaml** 文件(没有则新建一个)里面添加如下内容
```yaml
plugins:
  database:  # 插件类型
    neo4j: # 插件名
      default:                # 连接名称：default，可以是其他名字
        target: "bolt://192.168.0.93:7687"
        user: neo4j
        pwd: pwd123
      http:                # 连接名称
        target: "neo4j://127.0.0.1:7687"
        user: neo4j
        pwd: pwd123
```

3. 在需要使用的地方，直接使用
```go
使用默认 default 的client, 直接如下
neo4j.Cli.Session.Run("match (n:Person) return n", nil)

使用指定MySQL连接, 如下:
neo4j.GetCli("http").Session.Run("match (n:Person) return n", nil)
```