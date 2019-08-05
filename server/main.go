package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Received...")
		ioutil.ReadAll(r.Body)
		w.Write([]byte("Hello"))
		fmt.Println("Complete...")
	})
	http.ListenAndServe(":9000", nil)
}
