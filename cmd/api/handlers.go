package main

import (
	"fmt"
	"log"
	"net/http"
)

func (b *Broker) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	f, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	messageChan := make(chan string)

	b.newClients <- messageChan

	go func() {
		<-r.Context().Done()
		b.defunctClients <- messageChan
	}()

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Transfer-Encoding", "chunked")

	for {

		msg, open := <-messageChan

		if !open {
			break
		}

		_, err := fmt.Fprintf(w, "data: %s\n\n", msg)

		if err != nil {
			log.Println("Error occured while writing response", err)
		}
		f.Flush()
	}
}
