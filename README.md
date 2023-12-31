## 基于CloudWeGo的API网关开发

### 环境准备(以Debian/Ubuntu为例)

1. curl:

```
sudo apt update
sudo apt install curl
```

2. etcd:

下载：https://github.com/etcd-io/etcd/releases

解压得到 etcd （服务程序）和 etcdctl（命令行工具）

将这两个文件复制到 /usr/local/bin

```
curl -LO https://github.com/etcd-io/etcd/releases/download/v3.5.9/etcd-v3.5.9-linux-amd64.tar.gz

tar xvf etcd-v3.5.9-linux-amd64.tar.gz

cd etcd-v3.5.9-linux-amd64.tar.gz

sudo cp etcd /usr/local/bin/
sudo cp etcdctl /usr/local/bin/
```

3. golang:

```
curl -LO https://golang.google.cn/dl/go1.20.5.linux-amd64.tar.gz

sudo tar -C /usr/local -xzf go1.20.5.linux-amd64.tar.gz
```

​	sudo vim打开/etc/profile文件修改环境变量

```
sudo vim /etc/profile
```

​	追加一行

```
export PATH=$PATH:/usr/local/go/bin
```

​	使生效

```
source /etc/profile
```

​	验证环境

```
 go version
```

### 快速开始

1. 新建一个文件夹，将kitex-project与hertz-project放在这个文件夹下，文件结构如下

```

├── hertz-project
└── kitex-project
```

2. 开启etcd负载均衡：新建一个终端，输入以下指令

```
etcd --log-level debug
```

3. 启动hertz-project：新建一个终端，进入/hertz-project文件夹下，输入以下指令

```
go build
./hertz-project
```

4. 启动kitex-project服务端：新建一个终端，进入/kitex-project文件夹下，输入以下指令

```
go run .
```

5. 使用curl对服务进行访问：新建一个终端，输入curl指令访问

```
curl -H "Content-Type: application/json" -X POST http://127.0.0.1:8888/add-student-info -d '{"id": 100, "name":"Emma","sex":"female", "college": {"name": "software college", "address": "逸夫"}, "email": ["emma@nju.com"]}'
```

```
 curl -H "Content-Type: application/json" -X GET http://127.0.0.1:8888/query?id=100
```

### API网关核心功能展示

#### 正确响应 HTTP POST请求，请求体为JSON格式

如“快速开始”所示

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
   
   	root := r.Group("/", rootMw()...)
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
   
   	rg.POST("/add-student-info",demo.Register)
   
   	rg.GET("/query",demo.Query)
   
   }
   ```


**ps：本项目是参考第一种方式实现的**

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

参考项目代码



### 性能测试和优化报告

#### 测试方法说明


+ 测试文件: 在 test/main_test.go 和 test/mainMux_test.go 两个文件中编写测试代码以分别性能测试和并发量测试
+ 测试流程: 首先构造数据, 然后发送 GET、POST 请求并获取响应、检查状态码
+ 测试工具: 使用 benchmark 和 apache benchmark 等工具进行性能测试


#### 性能测试数据
1. 使用benchmark,在指令 "go test -bench=. main_test.go" 下的测试结果
    ```
    goos: linux
    goarch: amd64
    cpu: 11th Gen Intel(R) Core(TM) i7-11800H @ 2.30GHz
    BenchmarkAddStudentInfo-16                  2475            492007 ns/op
    BenchmarkQueryStudentInfo-16                2493            431299 ns/op
    PASS
    ok      command-line-arguments  4.062s
    ```
    
2. 使用 Apache benchmark 在指令 "ab -n 1000 -c 100 http://127.0.0.1:8888/add-student-info" 下的 测试结果:
    ```
    This is ApacheBench, Version 2.3 <$Revision: 1843412 $>
    Copyright 1996 Adam Twiss, Zeus Technology Ltd, http://www.zeustech.net/
    Licensed to The Apache Software Foundation, http://www.apache.org/
    
    Benchmarking 127.0.0.1 (be patient)
    Completed 100 requests
    Completed 200 requests
    Completed 300 requests
    Completed 400 requests
    Completed 500 requests
    Completed 600 requests
    Completed 700 requests
    Completed 800 requests
    Completed 900 requests
    Completed 1000 requests
    Finished 1000 requests
    
    
    Server Software:        hertz
    Server Hostname:        127.0.0.1
    Server Port:            8888
    
    Document Path:          /add-student-info
    Document Length:        18 bytes
    
    Concurrency Level:      100
    Time taken for tests:   0.039 seconds
    Complete requests:      1000
    Failed requests:        0
    Non-2xx responses:      1000
    Total transferred:      170000 bytes
    HTML transferred:       18000 bytes
    Requests per second:    25568.25 [#/sec] (mean)
    Time per request:       3.911 [ms] (mean)
    Time per request:       0.039 [ms] (mean, across all concurrent requests)
    Transfer rate:          4244.73 [Kbytes/sec] received
    
    Connection Times (ms)
                  min  mean[+/-sd] median   max
    Connect:        0    2   0.7      2       3
    Processing:     0    2   0.6      2       4
    Waiting:        0    1   0.6      1       3
    Total:          3    4   0.4      4       5
    
    Percentage of the requests served within a certain time (ms)
      50%      4
      66%      4
      75%      4
      80%      4
      90%      4
      95%      4
      98%      4
      99%      4
     100%      5 (longest request)
    ```
    
3. pprof测试结果
   ![pprof.png](https://box.nju.edu.cn/f/a4f7d97fefd24fb5ac10/?dl=1)
4. flame graph
   ![flameGraph.png](https://box.nju.edu.cn/f/7c477eaa091c474e897c/?dl=1)


#### 优化方法说明

对hertz项目中的Query请求添加cache：

```
//  hertz-project/biz/handler/hertz/demo/student_service.go

func Query(ctx context.Context, c *app.RequestContext) {
	var err error
	var req demo.QueryReq

	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	//counter>=3时才进行正常访问流程，counter<3时都从cache获取

	if counter<3 {
		resp, ok := respCache[int(req.ID)]
		if ok {
			counter++
			c.JSON(consts.StatusOK, resp)
			return
		}
	}
	counter=0

......

	/*****************处理结果***************/
	c.JSON(consts.StatusOK, resp)
	respCache[int(req.ID)]=resp
}

```

通过添加对Query的Cache，减少对9999端口服务端的访问，从而减少访问时间

#### 优化后性能数据

指令 "go test -bench=. main_test.go" 下的**优化后**测试结果

```
goos: linux
goarch: amd64
cpu: 11th Gen Intel(R) Core(TM) i7-11800H @ 2.30GHz
BenchmarkAddStudentInfo-16                  2475            492007 ns/op
BenchmarkQueryStudentInfo-16                2493            431299 ns/op
PASS
ok      command-line-arguments  4.062s
```

对比**优化前**：

```
goos: linux
goarch: amd64
cpu: 11th Gen Intel(R) Core(TM) i7-11800H @ 2.30GHz
BenchmarkAddStudentInfo-16                  2778            454898 ns/op
BenchmarkQueryStudentInfo-16                3937            305655 ns/op
PASS
ok      command-line-arguments  4.285s
```

观察两组数据的“BenchmarkQueryStudentInfo-16”项，平均访问时间从 **431299ns/op**降低至**305655 ns/op**，可见优化是成功的