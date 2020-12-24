# hua
hua - a scheme as code web api framework

## 特性
* 协议即代码(go)
* 基于反射的rpc服务（go）
* 生成多语言客户端和服务端代码
* mock框架
* https双向认证
* 生成文档工具
* 基于websocket支持订阅模式接口

### Go1 API

#### Server

```go
package main

import (
	"github.com/duanqy/hua/example/api"
	"github.com/duanqy/hua/pkg/hua"
)

service := &example.UserService{}
service.GetUser = func(arg *api.GetUserArg) (*api.GetUserReply,error) {
	panic("not implemented")
}
http.ListenAndServe("127.0.0.1", hua.BuildServer(service))
```


#### Client
```go
package main

import (
	"github.com/duanqy/hua/v2/example/api"
	"github.com/duanqy/hua/v2/pkg/hua"
)

client := &api.UserService{}
hua.BuildClient(&client)
reply,err := client.GetUser(&api.GetUserArg{Account:"sam"})
```

### Mock

#### Mock Server
```go
package main

import (
	"github.com/duanqy/hua/example/api"
	"github.com/duanqy/hua/pkg/hua"
)

service := &example.UserService{}
huamock.Stub(service)
http.ListenAndServe("127.0.0.1", hua.BuildServer(service))
```
#### Mock Client
```go
package main

import (
	"github.com/duanqy/hua/v2/example/api"
	"github.com/duanqy/hua/v2/pkg/hua"
)

client := &api.UserService{}
huamock.Stub(&client)
reply,err := client.GetUser(&api.GetUserArg{Account:"sam"})
```

### Go2 API

#### Server

```go
package main

import (
	"github.com/duanqy/hua/v2/example/api"
	"github.com/duanqy/hua/v2/pkg/hua"
)

service := hua.NewService[example.UserService]()
service.GetUser = func(arg *api.GetUserArg) (*api.GetUserReply,error) {
	panic("not implemented")
}
http.ListenAndServe("xxx", service)
```


#### Client
```go
package main

import (
	"github.com/duanqy/hua/v2/example/api"
	"github.com/duanqy/hua/v2/pkg/hua"
)

client := hua.Dial[example.UserService]("http://")
reply,err := client.GetUser(&api.GetUserArg{Account:"sam"})
```
