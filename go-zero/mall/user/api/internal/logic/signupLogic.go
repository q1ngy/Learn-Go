// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"time"

	"github.com/q1ngy/Learn-Go/mall/user/api/internal/svc"
	"github.com/q1ngy/Learn-Go/mall/user/api/internal/types"
	"github.com/q1ngy/Learn-Go/mall/user/model"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var secret = []byte("secret")

type SignupLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSignupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SignupLogic {
	return &SignupLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SignupLogic) Signup(req *types.SignupRequest) (resp *types.SignupResponse, err error) {
	if req.RePassword != req.Password {
		return nil, errors.New("两次输入密码不一致")
	}

	logx.Infov(req) // json.Marshal(req)
	logx.Infof("req: %#v", req)

	// 0.判断用户是否已注册
	u, err := l.svcCtx.UserModel.FindOneByUsername(l.ctx, req.Username)
	if err != nil && !errors.Is(err, sqlx.ErrNotFound) {
		logx.Errorw(
			"内部错误",
			logx.Field("err", err),
		)
		return nil, errors.New("内部错误")
	}
	if u != nil {
		return nil, errors.New("用户名已存在")
	}
	// 1.加密
	h := md5.New()
	h.Write([]byte(req.Password)) // md5
	h.Write(secret)               // 加盐
	passwordStr := hex.EncodeToString(h.Sum(nil))

	user := &model.User{
		UserId:   time.Now().Unix(),
		Username: req.Username,
		Password: passwordStr,
		Gender:   int64(req.Gender),
	}
	_, err = l.svcCtx.UserModel.Insert(context.Background(), user)
	if err != nil {
		return nil, err
	}
	return &types.SignupResponse{Message: "success"}, nil
}
