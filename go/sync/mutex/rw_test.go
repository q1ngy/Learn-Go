package mutex

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

var wg *sync.WaitGroup
var rwMutex *sync.RWMutex

func TestRWMutex(t *testing.T) {
	wg = new(sync.WaitGroup)
	rwMutex = new(sync.RWMutex)

	//wg.Add(2)
	//go readData(1)
	//go readData(2)

	wg.Add(3)
	go writeData(1)
	go writeData(2)
	go readData(3)

	wg.Wait()
}

func writeData(i int) {
	defer wg.Done()

	fmt.Println(i, " 准备上锁")
	rwMutex.Lock()
	fmt.Println(i, " 正在写")
	time.Sleep(3 * time.Second)
	rwMutex.Unlock()
	fmt.Println(i, " 写完了")
}

func readData(i int) {
	defer wg.Done()

	fmt.Println(i, " 准备上锁")
	rwMutex.RLock()
	fmt.Println(i, " 正在读")
	time.Sleep(3 * time.Second)
	rwMutex.RUnlock()
	fmt.Println(i, " 读完了")

}
