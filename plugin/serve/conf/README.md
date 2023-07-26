### conf
1. 在入口，比如 main.go 里面隐式导入log包路径
```go 
import _ "github.com/fankane/go-utils/plugin/serve/conf"
```

2. 在运行文件根目录下的 **system_plugin.yaml** 文件(没有则新建一个)里面添加如下内容
```yaml
plugins:
  serve:  # 插件类型: 服务类
    conf: # 插件名
      conf_file: test_conf.yaml
      watch_change: true           # 监听文件更新, 默认false
      change_cron: @every 10s      # watch_change=true 时生效，为空时则实时更新
```
[cron语法](../../../utime/README.md)

3. 在需要使用的地方，直接使用
- 测试文件样例
```yaml
a: 12
b: hello
```
- 使用代码
```go
import "github.com/fankane/go-utils/plugin/serve/conf"
type AB struct {
    A int    `yaml:"a"`
    B string `yaml:"b"`
}

x := &AB{}
conf.Unmarshal(x)
}
```



