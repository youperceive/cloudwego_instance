package main

import (
	captcha "CloudWeGoInstance/VerifyCodeService/kitex_gen/captcha/captchaservice"
	"log"
	"net"

	"github.com/cloudwego/kitex/server"
)

func main() {
	addr, err := net.ResolveTCPAddr("tcp", "0.0.0.0:8000")
	if err != nil {
		log.Fatal(err)
		return
	}

	svr := captcha.NewServer(new(CaptchaServiceImpl), server.WithServiceAddr(addr))

	err = svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
