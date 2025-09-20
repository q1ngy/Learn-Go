package sync_test

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

var num = 10
var wg sync.WaitGroup
var mutex sync.Mutex

func TestSaleTickets(t *testing.T) {
	wg.Add(4)

	go saleTickets(1)
	go saleTickets(2)
	go saleTickets(3)
	go saleTickets(4)

	wg.Wait()
}

func saleTickets(name int) {
	defer wg.Done()
	for {
		mutex.Lock()
		if num > 0 {
			time.Sleep(time.Millisecond * 100)
			num--
			fmt.Println(name, " 售出一张票，剩余：", num)
		} else {
			mutex.Unlock()
			fmt.Println(name, " 售罄")
			break
		}
		mutex.Unlock()
	}
}
