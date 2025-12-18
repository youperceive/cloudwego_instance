package main

import (
	"log"
	"net"

	user_account "github.com/youperceive/cloudwego_instance/rpc/user_account/kitex_gen/user_account/useraccountservice"
	"github.com/youperceive/cloudwego_instance/rpc/verify_code/kitex_gen/verify_code/verifycodeservice"

	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/server"
)

func main() {
	addr, err := net.ResolveTCPAddr("tcp", "0.0.0.0:8001")
	if err != nil {
		panic("fail to link to tcp addr:" + err.Error())
	}

	userAccountServiceImpl := new(UserAccountServiceImpl)
	userAccountServiceImpl.VerifyCodeClient = verifycodeservice.MustNewClient("CaptchaService", client.WithHostPorts("0.0.0.0:8000"))

	svr := user_account.NewServer(
		userAccountServiceImpl,
		server.WithServiceAddr(addr),
	)

	err = svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
