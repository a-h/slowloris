package main

import (
	"bufio"
	"bytes"
	"flag"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"fmt"

	"github.com/fatih/color"
)

var delayBeforeFirstByte = flag.Duration("delayBeforeFirstByte", 0, "How long to wait before the first byte is served.")
var contentToServe = flag.String("content", `{ "response": "ok" }`, "The JSON content to serve.")
var delayBeforeLastByte = flag.Duration("delayBeforeLastByte", 0, "How long to wait before the last byte is served.")
var port = flag.Int("port", 8080, "The default port to listen on.")

func main() {
	// Parse flags.
	flag.Parse()
	log.Printf("Listening on port: %d", *port)
	log.Printf("Delay before sending first byte: %v", *delayBeforeFirstByte)
	log.Printf("Responding with content: %v", *contentToServe)
	log.Printf("Delay before sending last byte: %v", *delayBeforeLastByte)

	h := &SlowHandler{
		DelayBeforeFirstByte: *delayBeforeFirstByte,
		ContentToServe:       *contentToServe,
		DelayBeforeLastByte:  *delayBeforeLastByte,
	}

	// Setup the server.
	mux := http.NewServeMux()
	mux.Handle("/", h)
	err := http.ListenAndServe(":"+strconv.Itoa(*port), mux)

	if err != nil {
		log.Printf("Error received by server: %s", err.Error())
	}

	log.Printf("Stopped!")
}

// SlowHandler handles responses. Just slowly.
type SlowHandler struct {
	DelayBeforeFirstByte time.Duration
	ContentToServe       string
	DelayBeforeLastByte  time.Duration
	RequestID            int32
}

func (h *SlowHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	id := atomic.AddInt32(&h.RequestID, 1)

	logRequest(id, fmt.Sprintf("Received HTTP request from %s to %s", r.RemoteAddr, r.RequestURI))
	for k, v := range r.Header {
		logRequest(id, fmt.Sprintf("%s: %s", k, strings.Join(v, "")))
	}
	defer r.Body.Close()

	reader := bufio.NewReader(r.Body)

	var err error
	for {
		var line string
		line, err = readFullLine(reader)
		logRequest(id, line)
		if err != nil {
			break
		}
	}

	if err != nil && err != io.EOF {
		logRequest(id, fmt.Sprintf("Failed to read the body: %s", err))
	}

	logRequest(id, fmt.Sprintf("About to send reponse after delay of %v...", h.DelayBeforeFirstByte))

	time.Sleep(h.DelayBeforeFirstByte)

	logRequest(id, fmt.Sprintf("Sending first byte of %d.", len(h.ContentToServe)))
	_, err = w.Write([]byte{h.ContentToServe[0]})

	attemptFlush(id, w)

	if err != nil {
		logRequest(id, fmt.Sprintf("Error sending byte: %s", err.Error()))
	}

	logRequest(id, fmt.Sprintf("Waiting for %v before sending %d remaining bytes.", h.DelayBeforeLastByte, len(h.ContentToServe)-1))

	time.Sleep(h.DelayBeforeLastByte)

	logRequest(id, fmt.Sprintf("Sending %d remaining bytes.", len(h.ContentToServe)-1))
	_, err = w.Write([]byte(h.ContentToServe[1:]))

	if err != nil {
		logRequest(id, fmt.Sprintf("Error sending byte: %s", err.Error()))
	}

	logRequest(id, "Complete.")
}

func attemptFlush(id int32, w http.ResponseWriter) {
	if f, ok := w.(http.Flusher); ok {
		f.Flush()
	} else {
		logRequest(id, "Unable to flush the HTTP response because the system doesn't support it.")
	}
}

func logRequest(id int32, msg string) {
	cmb := fmt.Sprintf("%d - %s %s", id, time.Now().Format("2006-01-02 15:04:05"), msg)

	if id%2 != 0 {
		color.Green(cmb)
		return
	}

	color.Yellow(cmb)
}

// http://stackoverflow.com/questions/8757389/reading-file-line-by-line-in-go
// https://play.golang.org/p/1rvfb3brME
func readFullLine(reader *bufio.Reader) (string, error) {
	var buffer bytes.Buffer

	var line []byte
	var isPrefix bool
	var err error
	for {
		line, isPrefix, err = reader.ReadLine()

		buffer.Write(line)

		// If we've reached the end of the line, stop reading.
		if !isPrefix {
			break
		}

		// If we're hit an error, stop reading.
		if err != nil {
			break
		}
	}

	return buffer.String(), err
}
