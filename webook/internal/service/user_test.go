package service

import (
	"context"
	"errors"
	"testing"

	"github.com/q1ngy/Learn-Go/webook/internal/domain"
	"github.com/q1ngy/Learn-Go/webook/internal/repository"
	repomocks "github.com/q1ngy/Learn-Go/webook/internal/repository/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
)

func TestPasswordEncrypt(t *testing.T) {
	password := []byte("abc123##")
	fromPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	t.Log(string(fromPassword))
	assert.NoError(t, err)
	err = bcrypt.CompareHashAndPassword(fromPassword, []byte("abc123##"))
	assert.NoError(t, err)
}

func TestUserServiceImpl_Login(t *testing.T) {
	testCases := []struct {
		name string
		mock func(ctrl *gomock.Controller) (repo repository.UserRepository)

		// 预期输入
		ctx      context.Context
		email    string
		password string

		// 预期输出
		wantUser domain.User
		wantErr  error
	}{
		{
			name: "登陆成功",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().
					FindByEmail(gomock.Any(), "123@123.com").
					Return(domain.User{
						Email:    "123@123.com",
						Password: "$2a$10$XKfM/wY3hU25l2qx/uPHEOmRHrNrwiRrVitW3MpUofxcpVuCcrD0S",
					}, nil)
				return repo
			},
			email:    "123@123.com",
			password: "abc123##",
			wantUser: domain.User{
				Email:    "123@123.com",
				Password: "$2a$10$XKfM/wY3hU25l2qx/uPHEOmRHrNrwiRrVitW3MpUofxcpVuCcrD0S",
			},
		},
		{
			name: "没有该用户",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().
					FindByEmail(gomock.Any(), "12@123.com").
					Return(domain.User{}, repository.ErrUserNotFound)
				return repo
			},
			email:    "12@123.com",
			password: "abc123##",
			wantUser: domain.User{},
			wantErr:  ErrInvalidUserOrPassword,
		},
		{
			name: "密码错误",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().
					FindByEmail(gomock.Any(), "123@123.com").
					Return(domain.User{
						Email:    "123@123.com",
						Password: "$2a$10$XKfM/wY3hU25l2qx/uPHEOmRHrNrwiRrVitW3MpUofxcpVuCcrD0S",
					}, nil)
				return repo
			},
			email:    "123@123.com",
			password: "abc123#",
			wantUser: domain.User{},
			wantErr:  ErrInvalidUserOrPassword,
		},
		{
			name: "db 错误",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().
					FindByEmail(gomock.Any(), "123@123.com").
					Return(domain.User{}, errors.New("db 错误"))
				return repo
			},
			email:    "123@123.com",
			password: "abc123##",
			wantUser: domain.User{},
			wantErr:  errors.New("db 错误"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := tc.mock(ctrl)
			userService := NewUserService(repo)

			user, err := userService.Login(tc.ctx, tc.email, tc.password)
			assert.Equal(t, tc.wantErr, err) // 先比较 err
			assert.Equal(t, tc.wantUser, user)
		})
	}
}
