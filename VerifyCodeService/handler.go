package main

import (
	"CloudWeGoInstance/VerifyCodeService/kitex_gen/base"
	captcha "CloudWeGoInstance/VerifyCodeService/kitex_gen/captcha"
	"CloudWeGoInstance/VerifyCodeService/pkg/redis"
	"CloudWeGoInstance/VerifyCodeService/pkg/util"
	"context"
	"log"
	"os"
	"strconv"
	"time"

	redis_v9 "github.com/redis/go-redis/v9"
)

// CaptchaServiceImpl implements the last service interface defined in the IDL.
type CaptchaServiceImpl struct{}

// GenerateCaptcha implements the CaptchaServiceImpl interface.
func (s *CaptchaServiceImpl) GenerateCaptcha(ctx context.Context, req *captcha.GenerateCaptchaRequest) (resp *captcha.GenerateCaptchaResponse, err error) {

	exist, err := redis.Exists(ctx, req.Target, req.Purpose)
	if err != nil {
		log.Println(err)
		resp = &captcha.GenerateCaptchaResponse{
			BaseResp: &base.BaseResponse{
				Code: base.Code_DB_ERR,
				Msg:  "Internal error",
			},
		}
		return
	}

	if exist {
		resp = &captcha.GenerateCaptchaResponse{
			BaseResp: &base.BaseResponse{
				Code: base.Code_INVALID_PARAM,
				Msg:  "Captcha already exists",
			},
		}
		return
	}

	code := util.GenerateCode()

	// here need to deliver the code to user's target (email or phone)

	err = redis.SetWithCount(ctx, req.Target, req.Purpose, code, time.Duration(req.ExpireSeconds)*time.Second, int(req.MaxValidateTimes))
	if err != nil {
		log.Println(err)
		resp = &captcha.GenerateCaptchaResponse{
			BaseResp: &base.BaseResponse{
				Code: base.Code_DB_ERR,
				Msg:  "Internal error",
			},
		}
		return
	}

	resp = &captcha.GenerateCaptchaResponse{
		BaseResp: &base.BaseResponse{
			Code: base.Code_SUCCESS,
			Msg:  "",
		},
	}

	if os.Getenv("PRINT_CAPTCHA") == "true" {
		log.Printf("Generated captcha for target %s, purpose %s: %s\n", req.Target, req.Purpose, code)
	}

	return
}

// ValidateCaptcha implements the CaptchaServiceImpl interface.
func (s *CaptchaServiceImpl) ValidateCaptcha(ctx context.Context, req *captcha.ValidateCaptchaRequest) (resp *captcha.ValidateCaptchaResponse, err error) {

	code, remain, err := redis.GetAndDecrementCount(ctx, req.Target, req.Purpose)
	if err != nil {
		log.Println(err)
		if err == redis_v9.Nil {
			resp = &captcha.ValidateCaptchaResponse{
				BaseResp: &base.BaseResponse{
					Code: base.Code_INVALID_PARAM,
					Msg:  "not exists the code",
				},
				Valid: false,
			}
		} else {
			resp = &captcha.ValidateCaptchaResponse{
				BaseResp: &base.BaseResponse{
					Code: base.Code_DB_ERR,
					Msg:  "Internal error",
				},
				Valid: false,
			}
		}
		return
	}

	if code != req.Captcha {
		msg := "wrong captcha, remain count: " + strconv.Itoa(remain)
		if remain == 0 {
			msg = "no remain count, has been deleted"
		}

		resp = &captcha.ValidateCaptchaResponse{
			BaseResp: &base.BaseResponse{
				Code: base.Code_INVALID_PARAM,
				Msg:  msg,
			},
			Valid: false,
		}
		return
	}

	err = redis.Delete(ctx, req.Target, req.Purpose)
	if err != nil {
		log.Println(err)
		resp = &captcha.ValidateCaptchaResponse{
			BaseResp: &base.BaseResponse{
				Code: base.Code_DB_ERR,
				Msg:  "Internal error",
			},
			Valid: false,
		}
		return
	}

	resp = &captcha.ValidateCaptchaResponse{
		BaseResp: &base.BaseResponse{
			Code: base.Code_SUCCESS,
			Msg:  "",
		},
		Valid: true,
	}

	return
}
