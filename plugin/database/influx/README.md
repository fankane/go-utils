### InfluxDB
1. 在入口，比如 main.go 里面隐式导入包路径
```go 
import _ "github.com/fankane/go-utils/plugin/database/influx"
```

2. 在运行文件根目录下的 **system_plugin.yaml** 文件(没有则新建一个)里面添加如下内容
```yaml
plugins:
  database:  # 插件类型
    influx: # 插件名
      default:                # 连接名称：default，可以是其他名字
        url: "http://localhost:8086"
        bucket: "hufan"
        org: "hfOrganization"
        write_async: false #异步写数据
        token: "nxQGnVkbsUooXoU9usr47Com5xQzvAaqkGxPL50NBAeGr1ApH8xGAcT2HtJvBOrnOigxVLYo9eGddXHfoEqTbQ=="
      hufan_influxDB:                # 连接名称：default，可以是其他名字
        url: "http://localhost:8086"
        bucket: "xxx"
        org: "yyy"
        write_async: false #异步写数据
        token: "nxQGnVkbsUooXoU9usr47Com5xQzvAaqkGxPL50NBAeGr1ApH8xGAcT2HtJvBOrnOigxVLYo9eGddXHfoEqTbQ=="
```

3. 在需要使用的地方，直接使用
```go
使用默认 default 的log, 直接如下
influx.Cli.WritePoint(context.Background(), point)

使用指定influxDB Client连接, 如下:
influx.GetClient("hufan_influxDB").WritePoint(context.Background(), point)

查询influxDB 里面的数据
influx.QueryClient().Query(context.Background(), query)
```