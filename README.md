# grpc_consul
consul+grpc demo

# 安装
下载：  
```
https://github.com/protocolbuffers/protobuf/releases
https://github.com/protocolbuffers/protobuf/releases/download/v32.0/protoc-32.0-win64.zip
```
解压后，把bin目录加到环境变量里面，比如路径： C:\protoc-32.0-win64\bin , 测试生效
```
protoc --version

提示:
C:\Users\xx>protoc --version
libprotoc 32.0
```

go 依赖包
```
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

# 生成
windwos 下 cmd 到 proto 目录下，然后执行：
```
protoc --proto_path=. --go_out=paths=source_relative:./payment --go-grpc_out=paths=source_relative:./payment payment.proto

```
ps: 需要先在 proto 目录下创建对应proto文件的子文件夹, 否则提示文件夹不存在


# 引用
新建 xxx 模块文件夹, 内部创建服务端, 需要go.mod文件, 在 xxx/go.mod 中引入 proto 模块
```
require (
	proto v0.0.1 // 这个会自动增加的
)

replace proto => ../proto

replace 用来告诉 Go：proto 模块实际上在本地 ../proto 目录
```

Tip：xx模块引用完可能需要执行 go mod tidy 一下

# consul
```
安装（windows, cmd msinfo32, 查看 系统类型 这一项）
https://developer.hashicorp.com/consul/install
```

启动开发模式
```
consul agent -dev
访问：http://127.0.0.1:8500/ui  (可以配置密码访问)
```