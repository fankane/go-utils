## plugin 下面的各个插件使用方法

### log
1. 在入口，比如 main.go 里面隐式导入log包路径
```go 
import _ "github.com/fankane/go-utils/plugin/log"
```

2. 在运行文件根目录下的 **system_plugin.yaml** 文件(没有则新建一个)里面添加如下内容
```yaml
plugins:
  log:  # 插件类型: log
    zap_sugar: # 插件名
        default:              # 日志名称：default，可以是其他名字
          level: debug        # 日志级别, 默认info [debug, info, warn, error, panic]
          enable_stdout: true # 是否开启日志同步到控制台, 默认false
          filename: ./test.log # 日志文件，不填默认使用 ./log.log
          max_size: 1         # 日志文件滚动日志的大小 单位 MB
          max_age: 2          # 最大日志保留天数
          max_backups: 7      # 最大日志备份数量
          compress: false     # 是否压缩，默认：false
        logname2:               # 日志名称：logname2，可以是其他名字
          level: warn           # 日志级别, 默认info [debug, info, warn, error, panic]
          enable_stdout: true   # 是否开启日志同步到控制台, 默认false
          filename: ./test2.log # 日志文件，不填默认使用 ./log.log
          max_size: 10          # 日志文件滚动日志的大小 单位 MB
          max_age: 7            # 最大日志保留天数
          max_backups: 5        # 最大日志备份数量
          compress: true        # 是否压缩，默认：false
```

3. 在需要使用的地方，直接使用
```go
使用默认 default 的log, 直接如下
log.Logger.Debug("xxx")

使用指定 log, 如下:
log.GetLogger("logname2").Error("xxx")
```