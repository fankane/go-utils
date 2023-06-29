# go-utils
实际项目开发时，常用的工具方法汇总。

> 如果你需要什么常见的功能，这里没有的，可提起pull request, 或者联系本人添加 <br>
> 邮箱: fanhu1116@qq.com 

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


## archive
> 压缩文件相关操作
- rar: rar文件的相关操作
  - UnRar: 解压 rar 文件
- zip: zip文件的相关操作
  - UnZip: 解压 zip 文件

## string
> 字符串操作

## time
> 日期相关功能

## slice
- contain
  - InInterfaces,InStrings, InInts 等，判定某个切片里面是否存在某个具体的值; 支持基础类型 int, string, float
- transform
  - SliToInterfaces: 将普通切片转换成interface切片，例如 []int -> []interface; 支持基础类型 int, string, float

## uuid

## random

