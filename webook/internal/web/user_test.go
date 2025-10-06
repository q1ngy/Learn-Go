package web

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/q1ngy/Learn-Go/webook/internal/domain"
	"github.com/q1ngy/Learn-Go/webook/internal/service"
	svcmocks "github.com/q1ngy/Learn-Go/webook/internal/service/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestHttp(t *testing.T) {
	req, err := http.NewRequest(http.MethodPost, "/users/signup", bytes.NewReader([]byte("body")))
	fmt.Println(req)
	assert.NoError(t, err)
	recorder := httptest.NewRecorder()
	assert.Equal(t, http.StatusOK, recorder.Code)
}

func TestMock(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	svc := svcmocks.NewMockUserService(ctrl)
	svc.EXPECT().SignUp(gomock.Any(), domain.User{
		Id:    1,
		Email: "123@qq.com",
	}).Return(errors.New("db err"))
	err := svc.SignUp(context.Background(), domain.User{
		Id:    1,
		Email: "123@qq.com",
	})
	t.Log(err)
}

func TestEmailPattern(t *testing.T) {
	testCases := []struct {
		name  string
		email string
		match bool
	}{
		{
			name:  "没 @",
			email: "123%qq.com",
			match: false,
		},
		{
			name:  "没后缀",
			email: "123@qq",
			match: false,
		},
		{
			name:  "合法邮箱",
			email: "123@qq.com",
			match: true,
		},
	}

	h := NewUserHandler(nil, nil)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			match, err := h.emailRexExp.MatchString(tc.email)
			require.NoError(t, err)
			require.Equal(t, match, tc.match)
		})
	}

}

func TestUserHandler_SignUp(t *testing.T) {
	testCases := []struct {
		name string
		mock func(ctrl *gomock.Controller) (service.UserService, service.CodeService)

		// 预期输入
		reqBuilder func(t *testing.T) *http.Request

		// 预期输出
		wantCode int
		wantBody string
	}{
		{
			name: "注册成功",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				userService := svcmocks.NewMockUserService(ctrl)
				userService.EXPECT().SignUp(gomock.Any(), domain.User{
					Email:    "123@qq.com",
					Password: "abc123##",
				}).Return(nil)
				codeService := svcmocks.NewMockCodeService(ctrl)
				return userService, codeService
			},
			reqBuilder: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(http.MethodPost, "/users/signup", bytes.NewReader([]byte(`{
					"email": "123@qq.com",
					"password": "abc123##",
					"confirmPassword": "abc123##"}`)))
				req.Header.Set("Content-Type", "application/json")
				assert.NoError(t, err)
				return req
			},
			wantCode: http.StatusOK,
			wantBody: "注册成功",
		},
		{
			name: "邮箱格式错误",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				userService := svcmocks.NewMockUserService(ctrl)
				codeService := svcmocks.NewMockCodeService(ctrl)
				return userService, codeService
			},
			reqBuilder: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(http.MethodPost, "/users/signup", bytes.NewReader([]byte(`{
					"email": "123@qq",
					"password": "abc123##",
					"confirmPassword": "abc123##"}`)))
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("Content-Type", "application/json")
				assert.NoError(t, err)
				return req
			},
			wantCode: http.StatusOK,
			wantBody: "非法邮箱格式",
		},
		{
			name: "两次输入密码不同",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				userService := svcmocks.NewMockUserService(ctrl)
				codeService := svcmocks.NewMockCodeService(ctrl)
				return userService, codeService
			},
			reqBuilder: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(http.MethodPost, "/users/signup", bytes.NewReader([]byte(`{
					"email": "123@qq.com",
					"password": "abc123##",
					"confirmPassword": "abc123###"}`)))
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("Content-Type", "application/json")
				assert.NoError(t, err)
				return req
			},
			wantCode: http.StatusOK,
			wantBody: "两次输入密码不对",
		},
		{
			name: "密码格式不对",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				userService := svcmocks.NewMockUserService(ctrl)
				codeService := svcmocks.NewMockCodeService(ctrl)
				return userService, codeService
			},
			reqBuilder: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(http.MethodPost, "/users/signup", bytes.NewReader([]byte(`{
					"email": "123@qq.com",
					"password": "abc123",
					"confirmPassword": "abc123"}`)))
				req.Header.Set("Content-Type", "application/json")
				assert.NoError(t, err)
				return req
			},
			wantCode: http.StatusOK,
			wantBody: "密码必须包含字母、数字、特殊字符，并且不少于八位",
		},
		{
			name: "邮箱冲突",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				userService := svcmocks.NewMockUserService(ctrl)
				userService.EXPECT().SignUp(gomock.Any(), domain.User{
					Email:    "123@123.com",
					Password: "abc123##",
				}).Return(service.EmailDuplicateErr)
				codeService := svcmocks.NewMockCodeService(ctrl)
				return userService, codeService
			},
			reqBuilder: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(http.MethodPost, "/users/signup", bytes.NewReader([]byte(`{
					"email": "123@123.com",
					"password": "abc123##",
					"confirmPassword": "abc123##"}`)))
				req.Header.Set("Content-Type", "application/json")
				assert.NoError(t, err)
				return req
			},
			wantCode: http.StatusOK,
			wantBody: "邮箱冲突",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userService, codeService := tc.mock(ctrl)
			userHandler := NewUserHandler(userService, codeService)

			server := gin.Default()
			userHandler.RegisterRoutes(server)

			req := tc.reqBuilder(t)
			recorder := httptest.NewRecorder()

			server.ServeHTTP(recorder, req)

			assert.Equal(t, tc.wantCode, http.StatusOK)
			assert.Equal(t, tc.wantBody, recorder.Body.String())
		})
	}
}
