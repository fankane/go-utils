### Elasticsearch
1. 在入口，比如 main.go 里面隐式导入mysql包路径
```go 
import _ "github.com/fankane/go-utils/plugin/database/es"
```

2. 在运行文件根目录下的 **system_plugin.yaml** 文件(没有则新建一个)里面添加如下内容
```yaml
plugins:
  database:  # 插件类型: 
    elasticsearch: # 插件名
      default:                # 连接名称：default，可以是其他名字
        addr: ["http://192.168.0.93:9200"]
        user:
        pwd:
      es2:                # 连接名称
        addr: ["http://127.0.0.1:9200"]
        user:
        pwd:
```

3. 在需要使用的地方，直接使用
```go
使用默认 default 的es客户端, 直接如下
es.Cli

使用指定MySQL连接, 如下:
es.GetClient("es2").Query()
```