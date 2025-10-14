package cache

import (
	"context"
	"fmt"
	"testing"

	"github.com/q1ngy/Learn-Go/webook/internal/repository/cache/redismocks"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestRedisCodeCache_Set(t *testing.T) {
	keyFunc := func(biz, phone string) string {
		return fmt.Sprintf("phone_code:%s:%s", biz, phone)
	}
	testCases := []struct {
		name    string
		mock    func(ctrl *gomock.Controller) redis.Cmdable
		ctx     context.Context
		biz     string
		phone   string
		code    string
		wantErr error
	}{
		{
			name: "设置成功",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				res := redismocks.NewMockCmdable(ctrl)
				cmd := redis.NewCmd(context.Background())
				cmd.SetErr(nil)
				cmd.SetVal(int64(0))
				res.EXPECT().Eval(gomock.Any(), luaSetCode,
					[]string{keyFunc("test", "13888888888")},
					[]any{"123456"}).Return(cmd)
				return res
			},
			biz:     "test",
			phone:   "13888888888",
			code:    "123456",
			wantErr: nil,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			codeCache := NewRedisCodeCache(tc.mock(ctrl))
			err := codeCache.Set(context.Background(), tc.biz, tc.phone, tc.code)

			assert.Equal(t, tc.wantErr, err)
		})
	}
}
