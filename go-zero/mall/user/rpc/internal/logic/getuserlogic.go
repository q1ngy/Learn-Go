package logic

import (
	"context"

	"github.com/q1ngy/Learn-Go/mall/user/rpc/internal/svc"
	"github.com/q1ngy/Learn-Go/mall/user/rpc/types/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserLogic {
	return &GetUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserLogic) GetUser(in *user.GetUserReq) (*user.GetUserResp, error) {
	u, _ := l.svcCtx.UserModel.FindOneByUserId(l.ctx, in.UserId)
	return &user.GetUserResp{UserId: u.UserId, Username: u.Username, Gender: u.Gender}, nil
}
