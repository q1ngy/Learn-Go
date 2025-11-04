// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"

	"github.com/q1ngy/Learn-Go/mall/order/api/internal/svc"
	"github.com/q1ngy/Learn-Go/mall/order/api/internal/types"
	"github.com/q1ngy/Learn-Go/mall/user/rpc/types/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type SearchLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSearchLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchLogic {
	return &SearchLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SearchLogic) Search(req *types.SearchRequest) (resp *types.SearchResponse, err error) {
	u, err := l.svcCtx.UserRpc.GetUser(l.ctx, &user.GetUserReq{UserId: 1})
	if err != nil {
		logx.Errorw("order rpc call err", logx.Field("err", err))
	}
	return &types.SearchResponse{OrderID: "1", Status: 100, Username: u.GetUsername()}, nil
}
