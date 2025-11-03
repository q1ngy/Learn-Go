// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package handler

import (
	"net/http"

	"github.com/q1ngy/Learn-Go/mall/user/api/internal/logic"
	"github.com/q1ngy/Learn-Go/mall/user/api/internal/svc"
	"github.com/q1ngy/Learn-Go/mall/user/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func SignupHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.SignupRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewSignupLogic(r.Context(), svcCtx)
		resp, err := l.Signup(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
