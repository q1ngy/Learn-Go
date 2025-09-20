package ctx_test

import (
	"context"
	"errors"
	"fmt"
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	// 哨兵异常
	UserDefinedErr = errors.New("user defined err")
)

func TestCtxCancel(t *testing.T) {
	parent, pCancel := context.WithCancelCause(context.Background())
	defer pCancel(UserDefinedErr)
	for i := 0; i < 3; i++ {
		go func(ctx context.Context) {
			for {
				select {
				case <-ctx.Done():
					return
				default:
					fmt.Printf("%d parent is working\n", i)
					time.Sleep(time.Millisecond * 100)
				}
			}
		}(parent)
	}

	child, cCancel := context.WithCancelCause(parent)
	defer cCancel(UserDefinedErr)
	for i := 0; i < 5; i++ {
		go func(ctx context.Context) {
			for {
				select {
				case <-ctx.Done():
					return
				default:
					fmt.Printf("%d child is working\n", i)
					time.Sleep(time.Millisecond * 100)
				}
			}
		}(child)
	}
	assert.Equal(t, runtime.NumGoroutine(), 10)
	time.Sleep(time.Second)

	pCancel(UserDefinedErr)
	time.Sleep(time.Second)
	assert.Equal(t, runtime.NumGoroutine(), 2)
	assert.Equal(t, parent.Err(), context.Canceled)
	assert.Equal(t, child.Err(), context.Canceled)
	assert.Equal(t, context.Cause(parent), UserDefinedErr)
	assert.Equal(t, context.Cause(child), UserDefinedErr)
	assert.Equal(t, <-parent.Done(), struct{}{})
	assert.Equal(t, <-child.Done(), struct{}{})

}
