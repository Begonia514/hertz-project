// Code generated by hertz generator. DO NOT EDIT.

package router

import (
	"github.com/cloudwego/hertz/pkg/app/server"
	hello_example "hertz-project/biz/router/hello/example"
	hertz_demo "hertz-project/biz/router/hertz/demo"
)

// GeneratedRegister registers routers generated by IDL.
func GeneratedRegister(r *server.Hertz) {
	//INSERT_POINT: DO NOT DELETE THIS LINE!
	hertz_demo.Register(r)

	hello_example.Register(r)
}
