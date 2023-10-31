### nacos

[nacos介绍](https://help.aliyun.com/document_detail/130146.html)

1. 在入口，比如 main.go 里面隐式导入nacos包路径
```go 
import _ "github.com/fankane/go-utils/plugin/serve/nacos"
```

2. 在运行文件根目录下的 **system_plugin.yaml** 文件(没有则新建一个)里面添加如下内容
```yaml
plugins:
  serve:  # 插件类型: 服务类
    nacos: # 插件名
      default:
        end_point: "localhost:8848"
        namespace_id: "xxx"
        timeout_ms: 5000
        username: "xxx"
        password: "xxx"
        data_id: "xxx"
        group: "xxx"
```

3. 在需要使用的地方，直接使用

- 使用代码
```go
import "github.com/fankane/go-utils/plugin/serve/nacos"
type AB struct {
    A int    `toml:"a"`
    B string `toml:"b"`
}

x := &AB{}

nacos.Cli.ParseListenConfig(func(newData string, v interface{}) error {
    if _, dErr := toml.Decode(newData, v); dErr != nil {
        log.Println("Decode err:", dErr)
        return dErr
    }
    return nil
}, x)

```