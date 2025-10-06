package repository

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/q1ngy/Learn-Go/webook/internal/domain"
	"github.com/q1ngy/Learn-Go/webook/internal/repository/cache"
	cachemocks "github.com/q1ngy/Learn-Go/webook/internal/repository/cache/mocks"
	"github.com/q1ngy/Learn-Go/webook/internal/repository/dao"
	daomocks "github.com/q1ngy/Learn-Go/webook/internal/repository/dao/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestCachedUserRepository_FindById(t *testing.T) {
	nowMs := time.Now().UnixMilli()
	now := time.UnixMilli(nowMs) // 截取至毫秒
	testCases := []struct {
		name string
		mock func(ctrl *gomock.Controller) (dao.UserDao, cache.UserCache)

		// 预期输入
		ctx context.Context
		uid int64

		// 预期输出
		wantUser domain.User
		wantErr  error
	}{
		{
			name: "查找用户，未命中缓存",
			mock: func(ctrl *gomock.Controller) (dao.UserDao, cache.UserCache) {
				uid := int64(123)
				userDao := daomocks.NewMockUserDao(ctrl)
				userCache := cachemocks.NewMockUserCache(ctrl)
				userCache.EXPECT().Get(gomock.Any(), uid).Return(domain.User{}, cache.ErrKeyNotExist)
				userDao.EXPECT().FindById(gomock.Any(), uid).Return(dao.User{
					Id:       123,
					Nickname: "nickname",
					Birthday: nowMs,
					Email: sql.NullString{
						String: "123@123.com",
						Valid:  true,
					},
					CTime: nowMs,
				}, nil)
				userCache.EXPECT().Set(gomock.Any(), domain.User{
					Id:       123,
					Nickname: "nickname",
					Birthday: now,
					Email:    "123@123.com",
					CTime:    now,
				})
				return userDao, userCache
			},
			uid: 123,
			wantUser: domain.User{
				Id:       123,
				Nickname: "nickname",
				Birthday: now,
				Email:    "123@123.com",
				CTime:    now,
			},
			wantErr: nil,
		},
		{
			name: "查找用户，命中缓存",
			mock: func(ctrl *gomock.Controller) (dao.UserDao, cache.UserCache) {
				uid := int64(123)
				userDao := daomocks.NewMockUserDao(ctrl)
				userCache := cachemocks.NewMockUserCache(ctrl)
				userCache.EXPECT().Get(gomock.Any(), uid).Return(domain.User{
					Id:       123,
					Nickname: "nickname",
					Birthday: now,
					Email:    "123@123.com",
					CTime:    now,
				}, nil)
				return userDao, userCache
			},
			uid: 123,
			wantUser: domain.User{
				Id:       123,
				Nickname: "nickname",
				Birthday: now,
				Email:    "123@123.com",
				CTime:    now,
			},
			wantErr: nil,
		},
		{
			name: "未找到用户",
			mock: func(ctrl *gomock.Controller) (dao.UserDao, cache.UserCache) {
				uid := int64(123)
				userDao := daomocks.NewMockUserDao(ctrl)
				userCache := cachemocks.NewMockUserCache(ctrl)
				userCache.EXPECT().Get(gomock.Any(), uid).Return(domain.User{}, cache.ErrKeyNotExist)
				userDao.EXPECT().FindById(gomock.Any(), uid).Return(dao.User{}, dao.ErrRecordNotFound)
				return userDao, userCache
			},
			uid:      123,
			wantUser: domain.User{},
			wantErr:  dao.ErrRecordNotFound,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userDao, userCache := tc.mock(ctrl)
			repo := NewCachedUserRepository(userDao, userCache)

			user, err := repo.FindById(context.Background(), tc.uid)
			assert.Equal(t, err, tc.wantErr)
			assert.Equal(t, user, tc.wantUser)
		})
	}
}
