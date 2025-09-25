package serivce

import (
	"fmt"
	"math/rand/v2"
	"testing"
)

func TestGenerateCode(t *testing.T) {
	// rand.Intn(n) 会返回区间 [0, n) 内的整数，不包括 n。
	code := rand.IntN(1000000)

	t.Log(code)
	t.Log(fmt.Sprintf("%06d", code))
	t.Log(fmt.Sprintf("%06d", 1)) // 不足补0，超过不截
}
