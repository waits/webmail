// Command webmail runs a webmail server over HTTPS.
package main

import (
	// "fmt"
	"log"
	"net/http"

	"github.com/waits/webmail/handler"
)

const addr = ":8080"

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handler.Index)
	mux.HandleFunc("/compose", handler.Compose)
	mux.HandleFunc("/mail/", handler.Message)
	mux.HandleFunc("/send", handler.Send)
	mux.HandleFunc("/static/style.css", handler.Static)

	s := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	log.Printf("Listening on %s\n", addr)
	log.Fatal(s.ListenAndServe())
}
