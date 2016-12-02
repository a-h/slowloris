package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
)

var url = flag.String("url", "http://localhost:8080", "The URL to read from.")

func main() {
	flag.Parse()

	resp, err := http.Get(*url)

	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		return
	}

	defer resp.Body.Close()
	buf := make([]byte, 1)

	received := 0
	for {
		read, err := resp.Body.Read(buf)

		if err != nil && err != io.EOF {
			fmt.Printf("Failed to read the body: %s\n", err.Error())
			break
		}

		if read < 1 {
			fmt.Print("\n")
			break
		}

		received++
		fmt.Print(".")
	}

	fmt.Printf("%d bytes received.\n", received)
}
