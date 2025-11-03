// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/q1ngy/Learn-Go/mall/user/api/internal/svc"
	"github.com/q1ngy/Learn-Go/mall/user/api/internal/types"
	"github.com/zeromicro/go-zero/core/stores/sqlx"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func passwordMd5(password string) string {
	h := md5.New()
	h.Write([]byte(password))
	h.Write(secret)
	return hex.EncodeToString(h.Sum(nil))
}

func (l *LoginLogic) Login(req *types.LoginRequest) (resp *types.LoginResponse, err error) {
	u, err := l.svcCtx.UserModel.FindOneByUsername(l.ctx, req.Username)
	if errors.Is(err, sqlx.ErrNotFound) {
		return &types.LoginResponse{Message: "用户名不存在"}, nil
	}
	if err != nil {
		return nil, errors.New("内部错误")
	}
	if u.Password != passwordMd5(req.Password) {
		return &types.LoginResponse{Message: "用户名或密码错误"}, nil
	}

	now := time.Now().Unix()
	expire := l.svcCtx.Config.Auth.AccessExpire
	token, err := l.getJwtToken(l.svcCtx.Config.Auth.AccessSecret, now, expire, u.UserId)
	if err != nil {
		logx.Errorw("token err", logx.Field("err", err))
		return nil, err
	}

	return &types.LoginResponse{
		Message:      "登陆成功",
		AccessToken:  token,
		AccessExpire: strconv.FormatInt(now+expire, 10),
		Refreshafter: strconv.FormatInt(now+expire, 10),
	}, nil
}

// 生成JWT方法
func (l *LoginLogic) getJwtToken(secretKey string, iat, seconds, userId int64) (string, error) {
	claims := make(jwt.MapClaims)
	claims["exp"] = iat + seconds
	claims["iat"] = iat
	claims["userId"] = userId

	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claims

	return token.SignedString([]byte(secretKey))
}
