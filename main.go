// Command webmail runs a webmail server over HTTPS.
package main

import (
	"crypto/tls"
	"flag"
	"log"
	"net/http"

	"github.com/waits/webmail/handler"
	"github.com/waits/webmail/maildir"
	"golang.org/x/crypto/acme/autocert"
)

var (
	addr = flag.String("addr", ":8080", "address to listen on")
	auth = flag.String("auth", "tmp/imap.passwd", "path of passwd file")
	dir  = flag.String("maildir", "tmp/inbox", "path of maildir")
	host = flag.String("host", "", "server hostname")
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
	mux.HandleFunc("/static/skeleton.css", handler.Static)
	mux.HandleFunc("/static/webmail.css", handler.Static)
	mux.HandleFunc("/static/webmail.js", handler.Static)

	log.Printf("Listening on %s\n", *addr)
	s := &http.Server{
		Addr:    *addr,
		Handler: mux,
	}
	if *host != "" {
		m := autocert.Manager{
			Prompt:     autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist(*host),
			Cache:      autocert.DirCache("certs"),
		}

		s80 := &http.Server{
			Addr:    ":80",
			Handler: m.HTTPHandler(nil),
		}
		go s80.ListenAndServe()

		s.TLSConfig = &tls.Config{GetCertificate: m.GetCertificate}
		log.Fatal(s.ListenAndServeTLS("", ""))
	} else {
		log.Fatal(s.ListenAndServe())
	}
}
