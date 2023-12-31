// Code generated by hertz generator.

package demo

import (
	"context"
	"os"
	"time"

	demo "hertz-project/biz/model/hertz/demo"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/adaptor"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	kclient "github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/client/genericclient"
	"github.com/cloudwego/kitex/pkg/generic"
	"github.com/cloudwego/kitex/pkg/klog"
)

// var PG ProviderGetter=ProviderGetter{lastThriftTime: 0}
// Register .
// @router /add-student-info [POST]
func Register(ctx context.Context, c *app.RequestContext) {
	var err error
	var req demo.Student

	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

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
	resp, err := cli.GenericCall(ctx, "Register", customReq)

	if err != nil {
		panic(err)
	}

	/*****************处理结果***************/
	c.JSON(consts.StatusOK, resp)
}

// Query .
// @router /query [GET]
//存储Query请求返回数据的cache
var respCache map[int]interface{} = make(map[int]interface{})
var counter int

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

	/*****************实现转化为http请求***************/
	httpReq, err := adaptor.GetCompatRequest(c.GetRequest())
	if err != nil {
		panic("get http req failed")
	}

	/*****************实现转化为 custom 请求***************/
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
	respCache[int(req.ID)]=resp
}

var p *generic.ThriftContentProvider = nil
var serviceNameCache string

/*****************全局的泛化client并会随着idl更新而实时更新***************/
var cli genericclient.Client = nil

func InitGenericClient(serviceName string) {

	serviceNameCache = serviceName
	idlContent, err := os.ReadFile("../kitex-project/idl/student.thrift")
	if err != nil {
		panic(err)
	}

	p, err = generic.NewThriftContentProvider(string(idlContent), map[string]string{})
	if err != nil {
		panic(err)
	}

	/*****************实现client的http泛化***************/
	g, err := generic.HTTPThriftGeneric(p)
	if err != nil {
		panic(err)
	}

	cli, err = genericclient.NewClient(serviceNameCache, g,
		kclient.WithHostPorts("127.0.0.1:9999"), //kclient.WithResolver(r),
	)
	if err != nil {
		panic(err)
	}

	go func() {
		ticker := time.NewTicker(time.Second * 10)
		for range ticker.C {
			UpdateIdl()

			klog.Info("update idl and generic client")
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
