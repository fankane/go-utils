### MongoDB
1. 在入口，比如 main.go 里面隐式导入mongo包路径
```go 
import _ "github.com/fankane/go-utils/plugin/database/mongodb"
```

2. 在运行文件根目录下的 **system_plugin.yaml** 文件(没有则新建一个)里面添加如下内容
```yaml
plugins:
  database:  # 插件类型
    mongo: # 插件名
      default:                # 连接名称：default，可以是其他名字
        host: "localhost"
        port: 27017
        user: ""
        pwd: ""
        connect_timeout_ms:  # 连接超时时限，单位：毫秒，不填默认 10s 
```

3. 在需要使用的地方，直接使用
```go
Cli.Cli.Database("testDB").Collection("testCollection").InsertOne(context.Background(), bson.M{})
```