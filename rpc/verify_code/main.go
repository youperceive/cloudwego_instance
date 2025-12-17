package main

import (
	"log"
	"net"
	"os"

	"github.com/youperceive/cloudwego_instance/rpc/verify_code/kitex_gen/captcha/captchaservice"

	"github.com/cloudwego/kitex/pkg/registry"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	etcd "github.com/kitex-contrib/registry-etcd"
)

func main() {
	var r registry.Registry
	etcdAddr := os.Getenv("ETCD_ADDR")

	if etcdAddr != "" {
		reg, err := etcd.NewEtcdRegistry([]string{etcdAddr})
		if err != nil {
			log.Println("fail to link to etcd"+err.Error(), "just running on host")
			r = registry.NoopRegistry
		} else {
			r = reg
		}
	} else {
		r = registry.NoopRegistry
	}

	addr, err := net.ResolveTCPAddr("tcp", "0.0.0.0:8000")
	if err != nil {
		panic("fail to link to addr" + err.Error())
	}

	svr := captchaservice.NewServer(
		new(CaptchaServiceImpl),
		server.WithRegistry(r),
		server.WithServerBasicInfo(
			&rpcinfo.EndpointBasicInfo{
				ServiceName: "VerifyCodeService",
			},
		),
		server.WithServiceAddr(addr),
	)

	err = svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
