package pool_test

import (
	"bytes"
	"fmt"
	"sync"
	"testing"
)

var (
	bufPool   = sync.Pool{New: func() any { return new(bytes.Buffer) }}
	bytesPool = sync.Pool{New: func() any {
		bytes := make([]byte, 0, 32*1024)
		return &bytes
	}}
)

func TestPool(t *testing.T) {
	s := useBuffer()
	fmt.Println(s)

	useBytes()
}

func useBuffer() string {
	buf := bufPool.Get().(*bytes.Buffer)
	buf.Reset()
	defer bufPool.Put(buf)

	buf.WriteString("hello")
	buf.WriteString(" world!")
	return buf.String()
}

func useBytes() {
	bp := bytesPool.Get().(*[]byte)
	b := (*bp)[:0] // 归零长度
	fmt.Println(b)
}
