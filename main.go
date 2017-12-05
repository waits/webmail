// Command webmail runs a webmail server over HTTPS.
package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/waits/webmail/handler"
	"github.com/waits/webmail/maildir"
)

var (
	addr = flag.String("addr", ":8080", "server hostname")
	dir  = flag.String("maildir", "tmp/inbox", "directory to store certificates in")
)

func main() {
	flag.Parse()

	maildir.Watch(*dir)

	mux := http.NewServeMux()
	mux.HandleFunc("/", handler.Index)
	mux.HandleFunc("/compose", handler.Compose)
	mux.HandleFunc("/mail/", handler.Message)
	mux.HandleFunc("/send", handler.Send)
	mux.HandleFunc("/static/style.css", handler.Static)

	s := &http.Server{
		Addr:    *addr,
		Handler: mux,
	}

	log.Printf("Listening on %s\n", *addr)
	log.Fatal(s.ListenAndServe())
}
