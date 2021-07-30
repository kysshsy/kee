# web框架 day1

## web框架的原理

在golang中要编写一个web服务器，tcp/ip协议等已经被golang包装，http服务器可以使用官方库net/http的实现。所以本web框架基于net/http实现，在其上完成动态路由等功能。

net/http完成http协议的解析，并提供基础的静态路由功能。

## 熟悉net/http接口

### 快速启动http服务器

```go
package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/hello", helloHandler)
  http.ListenAndServe(":9999", nil)
}
// handler echoes r.URL.Path
func indexHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "URL.Path = %q\n", req.URL.Path)
}

// handler echoes r.URL.Header
func helloHandler(w http.ResponseWriter, req *http.Request) {
	for k, v := range req.Header {
		fmt.Fprintf(w, "Header[%q] = %q\n", k, v)
	}
}
```

运行如上所示代码，可以运行一个绑定在9999端口的http服务器。

### 静态路由

所谓静态路由是指请求的路径与handler对应关系是写死的，与之相对的是动态路由。如果一个web框架提供动态路由，那么它可以以pattern`/:name/:subject`的方式去路由（找到路径对应的handler），其中以`:`起头的不是实际路径，而是参数，web框架会在收到请求后根据url获得该参数。比如请求`school.com/lihua/math`,参数name则为lihua,subject为math，这时就可以通过这两个参数返回李华的成绩了。

net/http提供静态路由，可以调用`http.HandleFunc(path, handler)`注册静态路由。访问path时会搜索到对应handler。

### net/http提供的接口

```go
func Handle(pattern string, handler Handler)
func HandleFunc(pattern string, handler func(ResponseWriter, *Request))
//这两种接口都可以在默认的router中注册静态路由，区别是一个参数类型是Handler，一个是func(ResponseWriteer, *Request)。

type Handler interface {
    ServeHTTP(ResponseWriter, *Request)
}
```

net/http将请求的上下文分割成两个部分 ResponseWriter和Request，依据名字也可以知道，一个用于写入相应，一个是请求的信息。

其实这两种接口都注册在一个默认router上，执行`http.ListenAndServe`时第二个参数为nil，则使用默认router。**web框架的关键**则是替换这个默认router，代替net/http提供的功能。让我们仔细观察ListenAndServe的第二个参数。

```go
func ListenAndServe(addr string, handler Handler) error
```

这个Handler在前一个代码框已有介绍，在注册静态路由时，可以注册一个简单处理请求的函数，也可以是一个实现了ServeHTTP接口的复杂结构，所以web框架的第一步是编写一个Handler interface，放到ListenAndServe第二个参数上。

## 实现

👋先我们实现一个静态路由engine来代替默认的router。

### 实现Http.Handler接口

```go
type Handler interface {
		ServeHTTP(ResponseWriter, *Request)
}
```

实现Handler接口只需要实现ServeHTTP方法。

```go
type engine struct {
	 handlers map[string]Handler
}

func (e *engine) ServeHTTP(writer ResponseWriter, req *Request) {
  	key := req.Method + "-" + req.URL.Path

	if value, ok := e.router[key]; ok {
		value(respWriter, req)
	} else {
  	// balabla
  }
}
// 注册静态路由
func (e *engine) Handle(method, pattern string, handler http.Handler) {
  key := method + "-" + pattern
  // 添加到 key-value对中
  e.handlers[key] = handler
} 

```

## 总结

1. 了解一般web框架的原理，比如gin。基于net/http的http服务器实现web框架。
2. 熟悉net/http提供的接口是实现web框架的第一步。



------

极客🐰博客：https://geektutu.com/post/gee.html