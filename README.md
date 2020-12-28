# 🌺 hua 
hua - a code as protocol web api framework

## 特性
* 协议即代码(go)
* 基于反射的rpc服务（go）
* 生成多语言客户端和服务端代码
* mock框架
* https双向认证
* 生成文档工具
* 基于websocket支持订阅模式接口

### Go Example

#### Server

```go
package main

import (
	"github.com/duanqy/hua/example/api"
	"github.com/duanqy/hua/pkg/hua"
)

service := &api.CalcService{}
service.Add = func(arg *api.AddArg) (*api.AddReply,error) {
	panic("not implemented")
}
http.ListenAndServe("127.0.0.1", hua.NewServer().Register(service))
```


#### Client
```go
package main

import (
	"github.com/duanqy/hua/example/api"
	"github.com/duanqy/hua/pkg/hua"
)

client := &api.CalcService{}
hua.BuildClient(&client)
reply,err := client.Add(context.Background(),&api.AddArg{Left:1,Right:2})
```

### Mock

#### Mock Server
```go
package main

import (
	"github.com/duanqy/hua/example/api"
	"github.com/duanqy/hua/pkg/hua"
)

service := &api.CalcService{}
huamock.Stub(service)
http.ListenAndServe("127.0.0.1",hua.NewServer().Register(service))
```
#### Mock Client
```go
package main

import (
	"github.com/duanqy/hua/example/api"
	"github.com/duanqy/hua/pkg/hua"
)

client := &api.CalcService{}
huamock.Stub(&client)
reply,err := client.Add(context.Background(),&api.AddArg{Left:1,Right:2})
```