package middleware

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/zeromicro/go-zero/rest"
)

// 全局中间件
// 记录所有请求的响应信息

type bodyCopy struct {
	http.ResponseWriter
	body *bytes.Buffer // 小本本
}

func newBodyCopy(w http.ResponseWriter) *bodyCopy {
	return &bodyCopy{
		ResponseWriter: w,
		body:           bytes.NewBuffer([]byte{}),
	}
}

func (bc bodyCopy) Write(b []byte) (int, error) {
	// 先给小本本里写
	bc.body.Write(b)
	// 再给http响应里写
	return bc.ResponseWriter.Write(b)
}

func CopyResp(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		bc := newBodyCopy(w)
		next(bc, r)
		fmt.Printf("req: %v resp: %v\n", r.URL, bc.body.String())
	}
}

func MiddlewareWithAnotherService(b bool) rest.Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if b {
			}
			next(w, r)
		}
	}
}
