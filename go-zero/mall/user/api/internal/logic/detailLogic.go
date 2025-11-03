// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"
	"fmt"

	"github.com/q1ngy/Learn-Go/mall/user/api/internal/svc"
	"github.com/q1ngy/Learn-Go/mall/user/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DetailLogic {
	return &DetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DetailLogic) Detail(req *types.DetailRequest) (resp *types.DetailResponse, err error) {
	fmt.Printf("JWT userId: %v\n", l.ctx.Value("userId"))
	userId := req.UserId
	u, err := l.svcCtx.UserModel.FindOneByUserId(l.ctx, userId)
	return &types.DetailResponse{Username: u.Username, Gender: u.Gender}, nil
}
