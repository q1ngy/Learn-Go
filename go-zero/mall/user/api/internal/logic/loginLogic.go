// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"

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
	return &types.LoginResponse{Message: "登陆成功"}, nil
}
