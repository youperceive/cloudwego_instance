package main

import (
	"log"
	"net"
	"os"

	user_account "github.com/youperceive/cloudwego_instance/rpc/user_account/kitex_gen/user_account/useraccountservice"
	"github.com/youperceive/cloudwego_instance/rpc/verify_code/kitex_gen/verify_code/verifycodeservice"

	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/server"
)

func main() {
	addr, err := net.ResolveTCPAddr("tcp", "0.0.0.0:8001")
	if err != nil {
		klog.Fatal("Init stage: ", "fail to link to tcp addr:"+err.Error())
	}

	verifyCodeAddr := os.Getenv("VERIFY_CODE_SERVICE_ADDR")
	if verifyCodeAddr == "" {
		klog.Fatal("Init stage: ", "env $VERIFY_CODE_SERVICE_ADDR is empty.")
	}

	userAccountServiceImpl := new(UserAccountServiceImpl)
	cli, err := verifycodeservice.NewClient(
		"verify-code-service",
		client.WithHostPorts(verifyCodeAddr),
	)
	if err != nil {
		klog.Fatal("Init stage: ", "fail to link verify-code-service. "+err.Error())
	}
	userAccountServiceImpl.VerifyCodeClient = cli

	svr := user_account.NewServer(
		userAccountServiceImpl,
		server.WithServiceAddr(addr),
	)

	err = svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
