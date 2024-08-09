### SQLite
1. 在入口，比如 main.go 里面隐式导入sqlite包路径
```go 
import _ "github.com/fankane/go-utils/plugin/database/sqlite"
```

2. 在运行文件根目录下的 **system_plugin.yaml** 文件(没有则新建一个)里面添加如下内容
```yaml
plugins:
  database:  # 插件类型
    sqlite: # 插件名
      default:                # sqlite 连接名称：default，可以是其他名字
        db_file: "F:/shareDir/sqlite/hf_test2.db"
        mode:  # 选填, 模式 [ro, rw, rwc, memory]
        cache: # 选填, [shared,private]
      test2:                # sqlite 连接名称：default，可以是其他名字
        db_file: "F:/shareDir/sqlite/hf_test2.db"
```

3. 在需要使用的地方，直接使用
```go
使用默认 default 的log, 直接如下
sqlite.DB.Query()

使用指定MySQL连接, 如下:
sqlite.GetDB("test2").Query()
```

4. 其他说明
- 驱动底层使用的 cgo 语法，所以编译运行的时候，需要设置 CGO_ENABLED=1
- 运行的时候，需要有 gcc 环境