// Code generated by hertz generator.

package demo

import (
	"context"

	demo "project2/biz/model/hertz/demo"
	// kitexdemo "project2/kitex_gen/hertz/demo"
	// studentservice "project2/kitex_gen/hertz/demo/studentservice"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/adaptor"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	kclient "github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/client/genericclient"
	"github.com/cloudwego/kitex/pkg/generic"
	//"golang.org/x/text/cases"
)

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

	cli := initGenericClient()

	httpReq, err := adaptor.GetCompatRequest(c.GetRequest())
	if err!=nil{
		panic("get http req failed")
	}

	customReq, err := generic.FromHTTPRequest(httpReq)
	if err!=nil{
		panic("get custom req failed")
	}

	resp, err := cli.GenericCall(ctx,"Register",customReq)

	if err!=nil{
		panic("resp error")
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

	cli := initGenericClient()

	httpReq, err := adaptor.GetCompatRequest(c.GetRequest())
	if err!=nil{
		panic("get http req failed")
	}

	customReq, err := generic.FromHTTPRequest(httpReq)
	if err!=nil{
		panic("get custom req failed")
	}

	resp, err := cli.GenericCall(ctx,"Query",customReq)

	if err!=nil{
		panic("resp error")
	}

	// resp := new(demo.Student)

	c.JSON(consts.StatusOK, resp)
}


func initGenericClient() genericclient.Client{
	p,err:=generic.NewThriftFileProvider("./idl/student.thrift")
	if err!=nil{
		panic(err)
	}

	g,err:=generic.HTTPThriftGeneric(p)
	if err!=nil{
		panic(err)
	}

	var cli genericclient.Client
/*
	switch service{
	case "Register":
		cli,err = genericclient.NewClient("Register", g ,
		kclient.WithHostPorts("127.0.0.1:9999"),
		)
	case "Query":

	}*/

	cli,err = genericclient.NewClient("service", g ,
		kclient.WithHostPorts("127.0.0.1:9999"),
	)

	if err!=nil{
		panic(err)
	}

	return cli
}