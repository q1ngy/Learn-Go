package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

func main() {
	start := time.Now()
	c := make(chan string)
	for _, url := range os.Args[1:] {
		go fetch(url, c)
	}
	for range os.Args[1:] {
		fmt.Println(<-c)
	}
	fmt.Printf("%.2fs elapsed\n", time.Since(start).Seconds())
}

func fetch(url string, c chan<- string) {
	start := time.Now()
	resp, err := http.Get(url)
	if err != nil {
		c <- fmt.Sprint(err)
		return
	}
	nbytes, err := io.Copy(io.Discard, resp.Body)
	resp.Body.Close()
	if err != nil {
		c <- fmt.Sprintf("while reading %s: %v", url, err)
		return
	}
	secs := time.Since(start).Seconds()
	c <- fmt.Sprintf("%.2fs %7d %s", secs, nbytes, url)
}
