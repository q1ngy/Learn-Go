package web

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHttp(t *testing.T) {
	req, err := http.NewRequest(http.MethodPost, "/users/signup", bytes.NewReader([]byte("body")))
	fmt.Println(req)
	assert.NoError(t, err)
	recorder := httptest.NewRecorder()
	assert.Equal(t, http.StatusOK, recorder.Code)
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
