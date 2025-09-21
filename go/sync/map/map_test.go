package map_test

import (
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMap(t *testing.T) {
	var m sync.Map
	var wg sync.WaitGroup

	writer, reader := 5, 5
	total := 1000

	wg.Add(writer)
	for i := 0; i < writer; i++ {
		go func(i int) {
			defer wg.Done()
			for i := 0; i < total; i++ {
				key := "k-" + strconv.Itoa(i)
				m.Store(key, i)
			}
		}(i)
	}

	wg.Add(reader)
	for i := 0; i < writer; i++ {
		go func(i int) {
			defer wg.Done()
			key := "k-" + strconv.Itoa(i)
			if _, ok := m.Load(key); ok {
				m.LoadOrStore(key, i)
			}
		}(i)
	}

	wg.Wait()
	count := 0
	m.Range(func(key, value any) bool {
		count++
		return true
	})
	assert.Equal(t, total, count)
}

func TestNormalMap(t *testing.T) {
	m := make(map[string]int)
	var wg sync.WaitGroup

	writers, readers := 4, 4
	n := 1000

	wg.Add(writers)
	for w := 0; w < writers; w++ {
		go func() {
			defer wg.Done()
			for i := 0; i < n; i++ {
				m["k-"+strconv.Itoa(i)] = i // 无锁写
			}
		}()
	}

	wg.Add(readers)
	for r := 0; r < readers; r++ {
		go func() {
			defer wg.Done()
			for i := 0; i < n; i++ {
				_ = m["k-"+strconv.Itoa(i)] // 无锁读
			}
		}()
	}

	wg.Wait()
}
