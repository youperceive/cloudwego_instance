package main

import (
	"log"
	"net"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/server"
	order "github.com/youperceive/cloudwego_instance/rpc/order/kitex_gen/order/orderservice"

	_ "github.com/youperceive/cloudwego_instance/rpc/order/pkg/mongo"
)

func main() {

	addr, err := net.ResolveTCPAddr("tcp", "0.0.0.0:8002")
	if err != nil {
		klog.Fatal("监听 0.0.0.0:8002 地址失败，" + err.Error())
	}

	svr := order.NewServer(
		new(OrderServiceImpl),
		server.WithServiceAddr(addr),
	)

	err = svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
