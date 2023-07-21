### pprof
1. 在入口，比如 main.go 里面隐式导入log包路径
```go 
import _ "github.com/fankane/go-utils/plugin/serve/pprof"
```

2. 在运行文件根目录下的 **system_plugin.yaml** 文件(没有则新建一个)里面添加如下内容
```yaml
plugins:
  serve:  # 插件类型: log
    pprof: # 插件名
      addr: localhost:6060 #pprof 程序监听地址，跟本身服务监听端口切勿重复
```

3. 服务启动后，可使用如下命令查看pprof 数据
```go
go tool pprof http://localhost:6060/debug/pprof/heap // 获取堆的相关数据
go tool pprof http://localhost:6060/debug/pprof/profile // 获取30s内cpu的相关数据
go tool pprof http://localhost:6060/debug/pprof/block // 在你程序调用 runtime.SetBlockProfileRate ，查看goroutine阻塞的相关数据
go tool pprof http://localhost:6060/debug/pprof/mutex // 在你程序调用 runtime.SetMutexProfileFraction，查看谁占用mutex
go tool pprof http://localhost:6060/debug/pprof/goroutine //查看当前所有运行的 goroutines 堆栈跟踪
go tool pprof http://localhost:6060/debug/pprof/allocs // 会采样自程序启动所有对象的内存分配信息（包括已经被GC回收的内存）
```