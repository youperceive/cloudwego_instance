package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/cloudwego/kitex/pkg/klog"
	base "github.com/youperceive/cloudwego_instance/rpc/user_account/kitex_gen/base"
	user "github.com/youperceive/cloudwego_instance/rpc/user_account/kitex_gen/user"
	"github.com/youperceive/cloudwego_instance/rpc/user_account/pkg/dao"
	"github.com/youperceive/cloudwego_instance/rpc/user_account/pkg/hash"

	"github.com/youperceive/cloudwego_instance/rpc/verify_code/kitex_gen/verify_code"
	"github.com/youperceive/cloudwego_instance/rpc/verify_code/kitex_gen/verify_code/verifycodeservice"
)

// UserAccountServiceImpl implements the last service interface defined in the IDL.
type UserAccountServiceImpl struct {
	VerifyCodeClient verifycodeservice.Client
}

var (
	internalErrMsg = "Internal Error."
	successMsg     = "Success."
)

func validateRegisterReq(req *user.RegisterRequest) error {
	var msg []string
	if req.Target == "" || req.Captcha == "" || req.Password == "" {
		msg = append(msg, "target, captcha or password is empty.")
	}
	if _, err := req.TargetType.Value(); err != nil {
		msg = append(msg, "RegisterType invalid.")
	}
	if len(msg) > 0 {
		return fmt.Errorf("%s", strings.Join(msg, ". "))
	}
	return nil
}

// Register implements the UserAccountServiceImpl interface.
func (s *UserAccountServiceImpl) Register(ctx context.Context, req *user.RegisterRequest) (resp *user.RegisterResponse, err error) {
	klogErr := func(msg string) {
		klog.Error(
			"method", "Register",
			"message", msg,
			"target", req.Target,
		)
	}

	if err := validateRegisterReq(req); err != nil {
		klogErr("fail to validate req params." + err.Error())
		resp = &user.RegisterResponse{
			BaseResp: &base.BaseResponse{
				Code: base.Code_INVALID_PARAM,
				Msg:  err.Error(),
			},
			UserId: nil,
		}
		return resp, nil
	}

	captchaReq := &verify_code.ValidateCaptchaRequest{
		Proj:    "user-account-service",
		BizType: "login",
		Target:  req.Target,
		Captcha: req.Captcha,
	}

	captchaResp, err := s.VerifyCodeClient.ValidateCaptcha(ctx, captchaReq)
	if err != nil {
		klogErr("fail to call verifyCodeClient.ValidateCaptcha()" + err.Error())
		resp = &user.RegisterResponse{
			BaseResp: &base.BaseResponse{
				Code: base.Code_SERVICE_ERR,
				Msg:  internalErrMsg,
			},
			UserId: nil,
		}
		return
	}
	if !captchaResp.Valid {
		msg := "fail to validate captcha."
		if captchaResp.BaseResp != nil {
			msg += captchaResp.BaseResp.Msg
		}
		klogErr(msg)
		resp = &user.RegisterResponse{
			BaseResp: &base.BaseResponse{
				Code: base.Code_INVALID_PARAM,
				Msg:  "fail to validate captcha.",
			},
			UserId: nil,
		}
		return
	}

	hashedPassword, err := hash.BCryptHash(req.Password)
	if err != nil {
		klogErr("fail to hash password." + err.Error())
		resp = &user.RegisterResponse{
			BaseResp: &base.BaseResponse{
				Code: base.Code_SERVICE_ERR,
				Msg:  internalErrMsg,
			},
			UserId: nil,
		}
		return
	}

	daoUser := &dao.User{
		Password: hashedPassword,
	}
	if req.Username != nil {
		daoUser.Username = *req.Username
	}
	if req.TargetType == base.TargetType_Phone {
		daoUser.Phone = &req.Target
	} else {
		daoUser.Email = &req.Target
	}
	daoUser.RegisterType = int8(req.TargetType)

	userID, err := dao.CreateUser(daoUser)
	if err != nil {
		klogErr("fail to create user." + err.Error())
		resp = &user.RegisterResponse{
			BaseResp: &base.BaseResponse{
				Code: base.Code_DB_ERR,
				Msg:  internalErrMsg,
			},
			UserId: nil,
		}
		return
	}

	resp = &user.RegisterResponse{
		BaseResp: &base.BaseResponse{
			Code: base.Code_SUCCESS,
			Msg:  successMsg,
		},
		UserId: &userID,
	}

	return
}
