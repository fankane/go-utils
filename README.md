# go-utils
实际项目开发时，常用的工具方法汇总。

> 如果你需要什么常见的功能，这里没有的，可提起pull request, 或者联系本人添加 <br>
> 邮箱: fanhu1116@qq.com 

## plugin
#### <font style="color: red">不需要写代码</font> 将一些基础功能，通过配置文件的形式，在服务启动的时候，自动加载，需要的时候，直接使用； <br>

```go
import	"github.com/fankane/go-utils/plugin"

plugin.Load()
```

[各种插件使用方法](plugin/README.md)

## file
> 文件相关操作
- 读
  - [x] DirFiles: 目录下面的文件名列表
  - [x] Content: 读取文件内容
  - [x] FileExist: 文件是否存在
  - [x] DirExist: 目录是否存在
- 写
  - [x] Mkdir: 创建目录
  - [x] DeleteFiles: 删除文件
  - [x] DeleteDir: 删除目录
  - [x] DeleteDirFilesWithPref: 删除目录下指定前缀的文件
- xlsx: xlsx 文件处理
- xls: xls 文件的读[建议优先考虑xlsx]
- csv: csv 文件处理

## archive
> 压缩文件相关操作
- rar: rar文件的相关操作
  - UnRar: 解压 rar 文件
- zip: zip文件的相关操作
  - CreateZip: 创建 zip 文件
  - UnZip: 解压 zip 文件

## string
> 字符串操作
  - 字符串转换
  - uuid
## time
> 日期相关功能

## slice
- contain
  - InInterfaces,InStrings, InInts 等，判定某个切片里面是否存在某个具体的值; 支持基础类型 int, string, float
- transform
  - ToInterfaceSli: 将普通切片转换成interface切片，例如 []int -> []interface; 支持基础类型 int, string, float
- compare
  - StrSliContentEqual: 比较两个字符串切片内容是否相同，忽略顺序

## random

## err
> 自定义error结构

## http
> http 方法包装