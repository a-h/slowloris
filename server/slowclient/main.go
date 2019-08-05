package main

import (
	"flag"
	"fmt"
	"net/http"
	"sync"
	"time"
)

var url = flag.String("url", "http://localhost:9000", "The URL to read from.")
var n = flag.Int("n", 5000, "Number of connections to make.")

type slowReader struct {
	delay time.Duration
}

func (sr slowReader) Read(p []byte) (n int, err error) {
	p[0] = 0
	time.Sleep(sr.delay)
	return 1, nil
}

func main() {
	flag.Parse()
	sr := slowReader{
		delay: time.Minute,
	}

	var wg sync.WaitGroup
	wg.Add(*n)
	for i := 0; i < *n; i++ {
		go func() {
			r, err := http.NewRequest("POST", *url, sr)
			if err != nil {
				fmt.Printf("invalid request: %v\n", err)
			}
			_, err = http.DefaultClient.Do(r)
			if err != nil {
				fmt.Printf("request err: %v\n", err)
			}
			fmt.Println("Sent request...")
		}()
	}
	wg.Wait()
	fmt.Println("Done")
}
