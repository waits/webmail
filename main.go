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
	auth = flag.String("auth", "tmp/imap.passwd", "path of passwd file")
	dir  = flag.String("maildir", "tmp/inbox", "directory to store certificates in")
)

func main() {
	flag.Parse()

	handler.LoadPasswd(*auth)
	maildir.Watch(*dir)

	mux := http.NewServeMux()
	mux.HandleFunc("/", handler.WithAuth(handler.Index))
	mux.HandleFunc("/compose", handler.WithAuth(handler.Compose))
	mux.HandleFunc("/mail/", handler.WithAuth(handler.Message))
	mux.HandleFunc("/send", handler.WithAuth(handler.Send))
	mux.HandleFunc("/static/normalize.css", handler.Static)
	mux.HandleFunc("/static/sakura.css", handler.Static)
	mux.HandleFunc("/static/webmail.css", handler.Static)

	s := &http.Server{
		Addr:    *addr,
		Handler: mux,
	}

	log.Printf("Listening on %s\n", *addr)
	log.Fatal(s.ListenAndServe())
}
