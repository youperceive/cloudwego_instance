package main

import (
	order "github.com/youperceive/cloudwego_instance/rpc/order/kitex_gen/order/orderservice"
	"log"

	_ "github.com/youperceive/cloudwego_instance/rpc/order/pkg/mongo"
)

func main() {
	svr := order.NewServer(new(OrderServiceImpl))

	err := svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
