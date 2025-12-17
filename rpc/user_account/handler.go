package main

import (
	"context"
	"fmt"
	"log"
	"strings"

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
	if err := validateRegisterReq(req); err != nil {
		resp = &user.RegisterResponse{
			BaseResp: &base.BaseResponse{
				Code: base.Code_INVALID_PARAM,
				Msg:  err.Error(),
			},
			UserId: nil, // temporary
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
		log.Println(err)
		resp = &user.RegisterResponse{
			BaseResp: &base.BaseResponse{
				Code: base.Code_SERVICE_ERR,
				Msg:  "Internal error.",
			},
			UserId: nil,
		}
		return
	}
	if !captchaResp.Valid {
		msg := ""
		if captchaResp.BaseResp != nil {
			msg = captchaResp.BaseResp.Msg
		}
		resp = &user.RegisterResponse{
			BaseResp: &base.BaseResponse{
				Code: base.Code_INVALID_PARAM,
				Msg:  "fail to validate captcha." + msg,
			},
			UserId: nil,
		}
		return
	}

	daoUser := &dao.User{
		Password: hash.Hash(req.Password),
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
		log.Println(err)
		resp = &user.RegisterResponse{
			BaseResp: &base.BaseResponse{
				Code: base.Code_DB_ERR,
				Msg:  "fail to create user.",
			},
			UserId: nil,
		}
		return
	}

	resp = &user.RegisterResponse{
		BaseResp: &base.BaseResponse{
			Code: base.Code_SUCCESS,
			Msg:  "",
		},
		UserId: &userID,
	}

	return
}
