// Code generated by hertz generator.

package demo

import (
	"context"
	"os"
	"time"

	demo "project2/biz/model/hertz/demo"
	// kitexdemo "project2/kitex_gen/hertz/demo"
	// studentservice "project2/kitex_gen/hertz/demo/studentservice"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/adaptor"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	kclient "github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/client/genericclient"
	"github.com/cloudwego/kitex/pkg/generic"
	"github.com/cloudwego/kitex/pkg/klog"
	// etcd "github.com/kitex-contrib/registry-etcd"
	//"golang.org/x/text/cases"
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

	// httpReq, errhttp := adaptor.GetCompatRequest(c.GetRequest())
	// if errhttp != nil {
	// 	c.String(consts.StatusBadRequest, errhttp.Error())
	// }
	// cli, errhitex := studentservice.NewClient("register", client.WithHostPorts("127.0.0.1:9999"))
	// if errhitex != nil {
	// 	c.String(consts.StatusBadRequest, errhitex.Error())
	// }

	// RegiReq := &kitexdemo.Student{
	// 	Id:   req.ID,
	// 	Name: req.Name,
	// 	College: &kitexdemo.College{
	// 		Name:    req.College.Name,
	// 		Address: req.College.Address,
	// 	},
	// 	Email: req.Email,
	// }
	// RegiResp, respError := cli.Register(context.Background())

	// cli := initGenericClient("add-student-info")

	httpReq, err := adaptor.GetCompatRequest(c.GetRequest())
	if err != nil {
		panic("get http req failed")
	}

	customReq, err := generic.FromHTTPRequest(httpReq)
	if err != nil {
		panic("get custom req failed")
	}

	resp, err := cli.GenericCall(ctx, "Register", customReq)

	if err != nil {
		// panic("resp error")
		panic(err)
	}

	//resp := new(demo.RegisterResp)

	c.JSON(consts.StatusOK, resp)
}

// Query .
// @router /query [GET]
func Query(ctx context.Context, c *app.RequestContext) {
	var err error
	var req demo.QueryReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	// cli := initGenericClient("query")

	httpReq, err := adaptor.GetCompatRequest(c.GetRequest())
	if err != nil {
		panic("get http req failed")
	}

	customReq, err := generic.FromHTTPRequest(httpReq)
	if err != nil {
		panic("get custom req failed")
	}

	resp, err := cli.GenericCall(ctx, "Query", customReq)

	if err != nil {
		panic("resp error")
	}

	// resp := new(demo.Student)

	c.JSON(consts.StatusOK, resp)
}

var p *generic.ThriftContentProvider = nil
var cli genericclient.Client = nil

func InitGenericClient(serviceName string) {

	// r, err := etcd.NewEtcdResolver([]string{"127.0.0.1:2379"})
	// if err != nil {
	// 	panic(err)
	// }

	idlContent, err := os.ReadFile("../project3/idl/student.thrift")
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
	idlContent, err := os.ReadFile("../project3/idl/student.thrift")
	if err != nil {
		panic(err)
	}

	err = p.UpdateIDL(string(idlContent), nil)
	if err != nil {
		panic(err)
	}
}
