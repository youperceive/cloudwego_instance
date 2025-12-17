package main

import (
	"CloudWeGoInstance/UserAccountService/kitex_gen/captcha/captchaservice"
	user "CloudWeGoInstance/UserAccountService/kitex_gen/user/useraccountservice"
	"log"
	"net"

	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/server"
)

func main() {
	addr, err := net.ResolveTCPAddr("tcp", "0.0.0.0:8001")
	if err != nil {
		panic("fail to link to tcp addr:" + err.Error())
	}

	userAccountServiceImpl := new(UserAccountServiceImpl)
	userAccountServiceImpl.CaptchaClient = captchaservice.MustNewClient("CaptchaService", client.WithHostPorts("0.0.0.0:8000"))

	svr := user.NewServer(
		userAccountServiceImpl,
		server.WithServiceAddr(addr),
	)

	err = svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
