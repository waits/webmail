// Package handler holds HTTP request handlers
package handler

import (
	"log"
	"net/http"
	"path/filepath"

	"github.com/waits/webmail/template"
)

const staticBase = "static/"

func Index(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL.Path)
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	template.Render(w, "index.html", nil)
}

func Compose(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL.Path)
	template.Render(w, "compose.html", nil)
}

func Static(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, staticBase + filepath.Base(r.URL.Path))
}
