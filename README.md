## 基于CloudWeGo的API网关开发

### API网关核心功能展示

#### 正确响应 HTTP POST请求，请求体为JSON格式



#### 根据请求路由确认目标服务和方法

1. 我们可以通过thrift与hertz快速构建项目和注册与确认路由

   ```shell
   hz client --idl=hertz-project/idl/xxx.thrift
   ```

   会在./biz/router下找到对应文件：本实例的对应文件为./biz/router/hertz/demo/student.go

   ```
   //  hertz-project/biz/router/hertz/demo/student.go
   ......
   func Register(r *server.Hertz) {
   
   	root := r.Group("/v1", rootMw()...)// 此处的“/v1”为手动添加，原本为“/”
   	root.POST("/add-student-info", append(_registerMw(), demo.Register)...)
   	root.GET("/query", append(_queryMw(), demo.Query)...)
   }
   ```

   通过该方法注册路由

2. 还可以自己在./router.go内部手动添加注册代码

   ```
   ...
   import (
   	"github.com/cloudwego/hertz/pkg/app/server"
   	handler "hertz-project/biz/handler"
   	//demo "hertz-project/biz/handler/hertz/demo"
   )
   
   func customizedRegister(r *server.Hertz) {
   	r.GET("/ping", handler.Ping)
   
   
   
   	rg := r.Group("/v1")
   
   	rg.POST("/add-student-info",demo.Register)
   
   	rg.GET("/query",demo.Query)
   
   }
   ```
   
   



#### 网关内的 IDL 管理模块，可为构造 Kitex Client 提供 IDL

本项目的IDL管理模块在hertz-project/biz/handler/hertz/demo/student_sesrvice.go内部，为InitGernericClient函数

```
//   hertz-project/biz/handler/hertz/demo/student_sesrvice.go

...

var p *generic.ThriftContentProvider = nil
var cli genericclient.Client = nil

func InitGenericClient(serviceName string) {


	idlContent, err := os.ReadFile("../kitex-project/idl/student.thrift")
	if err != nil {
		panic(err)
	}

	p, err = generic.NewThriftContentProvider(string(idlContent), map[string]string{})
	if err != nil {
		panic(err)
	}

	g, err := generic.HTTPThriftGeneric(p)
	if err != nil {
		panic(err)
	}

	cli, err = genericclient.NewClient(serviceName, g,
		kclient.WithHostPorts("127.0.0.1:9999"), //kclient.WithResolver(r),
	)
	if err != nil {
		panic(err)
	}

	go func() {
		ticker := time.NewTicker(time.Second * 10)
		for range ticker.C {
			UpdateIdl()
			klog.Info("update idl")
		}
	}()
}

func UpdateIdl() {
	idlContent, err := os.ReadFile("../kitex-project/idl/student.thrift")
	if err != nil {
		panic(err)
	}

	err = p.UpdateIDL(string(idlContent), nil)
	if err != nil {
		panic(err)
	}
}
```

该模块提供动态IDL功能，每10s自动从kitex/idl/student.thrift获取新的IDL

该模块触发于hertz-project/main.go内部

```
//    hertz-project/main.go

...

func main() {
	go http.ListenAndServe("localhost:8080",nil)
	h := server.Default()

	register(h)
	demo.InitGenericClient("studentservice")
	h.Spin()
}
```





#### 构造 Kitex 泛化调用客户端、发起请求并处理影响结果

泛化调用实现于hertz-project/biz/handler/hertz/demo/student_sesrvice.go内部，

```
//    hertz-project/biz/handler/hertz/demo/student_sesrvice.go


func Query(ctx context.Context, c *app.RequestContext) {
	...
	
	/*****************实现转化为http请求***************/
	httpReq, err := adaptor.GetCompatRequest(c.GetRequest())
	if err != nil {
		panic("get http req failed")
	}
	/*****************实现转化为custom请求***************/
	customReq, err := generic.FromHTTPRequest(httpReq)
	if err != nil {
		panic("get custom req failed")
	}
	/*****************实现泛化client发起请求***************/
	resp, err := cli.GenericCall(ctx, "Query", customReq)

	if err != nil {
		panic("resp error")
	}
	/*****************处理结果***************/
	c.JSON(consts.StatusOK, resp)
}

var p *generic.ThriftContentProvider = nil
/*****************全局的泛化client并会随着idl更新而实时更新***************/
var cli genericclient.Client = nil

func InitGenericClient(serviceName string) {
	...
	
	/*****************实现client的http泛化***************/
	g, err := generic.HTTPThriftGeneric(p)
	if err != nil {
		panic(err)
	}

	cli, err = genericclient.NewClient(serviceName, g,
		kclient.WithHostPorts("127.0.0.1:9999"), //kclient.WithResolver(r),
	)
	if err != nil {
		panic(err)
	}

	...
}


```

这里以query函数为例，client通过InitGenericClient函数定时变动并在内部泛化，

#### 编码：代码可读性，模块划分合理性，单元测试覆盖率等





### 性能测试和优化报告

#### 测试方法说明



#### 性能测试数据



#### 优化方法说明



#### 优化后性能数据






