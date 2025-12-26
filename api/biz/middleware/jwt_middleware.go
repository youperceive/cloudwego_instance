package middleware

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/youperceive/cloudwego_instance/rpc/user_account/pkg/token"
)

const UserIDKey = "user_id"

var jwtWhitelist = map[string]bool{
	"/create": true,
}

func JWTMiddleware() app.HandlerFunc {
	return func(c context.Context, ctx *app.RequestContext) {
		hlog.Info(ctx.FullPath())

		path := ctx.FullPath()
		if !jwtWhitelist[path] {
			ctx.Next(c)
			return
		}

		tokenStr := ctx.GetHeader("user_token")
		if len(tokenStr) == 0 {
			ctx.JSON(consts.StatusUnauthorized, map[string]interface{}{
				"baseResp": map[string]interface{}{
					"code": -1,
					"msg":  "Token为空，请先登录",
				},
			})
			ctx.Abort()
			return
		}
		hlog.Info(string(tokenStr))

		claims, err := token.VerifyToken(string(tokenStr))
		if err != nil {
			hlog.Error("Token解析失败:", err)
			ctx.JSON(consts.StatusUnauthorized, map[string]interface{}{
				"baseResp": map[string]interface{}{
					"code": -1,
					"msg":  err.Error(),
				},
			})
			ctx.Abort()
			return
		}

		ctx.Set(UserIDKey, claims.UserID)
		hlog.Debug("JWT校验通过，用户ID:", claims.UserID)

		ctx.Next(c)
	}
}
