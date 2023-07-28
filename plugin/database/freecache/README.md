### freecache
1. 在入口，比如 main.go 里面隐式导入freecache包路径
```go 
import _ "github.com/fankane/go-utils/plugin/database/freecache"
```

2. 在运行文件根目录下的 **system_plugin.yaml** 文件(没有则新建一个)里面添加如下内容
```yaml
plugins:
  database:  # 插件类型:
    freecache: # 插件名
      default:                # 连接名称：default，可以是其他名字
        cache_size: 102400
      cache2:                # 连接名称：default，可以是其他名字
        cache_size: 102400
```

3. 在需要使用的地方，直接使用
```go
使用默认 default 的freecache, 直接如下
freecache.Cache.Get()

使用指定freecache连接, 如下:
freecache.GetCache("cache2").Get()
```