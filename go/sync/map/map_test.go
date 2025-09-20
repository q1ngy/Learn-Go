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
