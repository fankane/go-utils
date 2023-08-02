### Prometheus
1. 在入口，比如 main.go 里面隐式导入prometheus包路径
```go 
import _ "github.com/fankane/go-utils/plugin/monitor/prometheus"
```

2. 在运行文件根目录下的 **system_plugin.yaml** 文件(没有则新建一个)里面添加如下内容
```yaml
plugins:
  monitor:  # 插件类型:
    prometheus: # 插件名
      port: 7701  #Prometheus服务监听端口号，不要和服务本身端口重复

```

3. 效果展示
- 3.1 什么都不使用的时候，默认采集go服务的相关数据，结合grafana面板即可看到效果
  - [Go-metrics 面板配置样例](https://grafana.com/grafana/dashboards/10826-go-metrics/)
  - ![avatar](../../image/go-metrics-panel.png)

- 3.2 自定义数据上报
```go
    g := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "test_hf",
	}, []string{"label1", "label2"})
	prometheus.MustRegister(g)
	for i := 0; i < 100; i++ {
		g.WithLabelValues("val1", "val2").Set(1.0)
		g.WithLabelValues("val1", "val3").Set(2.0)
		g.WithLabelValues("val2", "val3").Set(3.0)
		time.Sleep(time.Millisecond * 50)
}
```
- ![avatar](../../image/go-custom.png)

4. [prometheus & grafana 安装教程](https://juejin.cn/post/7130391327413370887)