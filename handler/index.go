// Package handler implements HTTP request handlers.
package handler

import (
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/waits/webmail/maildir"
	"github.com/waits/webmail/template"
)

const staticBase = "static/"

// Index serves the home page.
func Index(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL.Path)
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	template.Render(w, "index.html", maildir.Messages)
}

// Compose serves the compose mail page.
func Compose(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL.Path)
	template.Render(w, "compose.html", nil)
}

// Message serves the message detail page.
func Message(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL.Path)
	paths := strings.Split(r.URL.Path, "/")

	id := paths[len(paths)-1]
	msg, ok := maildir.Messages[id]
	if !ok {
		http.NotFound(w, r)
		return
	}

	switch method(r) {
	case "GET":
		template.Render(w, "message.html", msg)
	case "DELETE":
		maildir.DeleteMessage(msg.ID) // FIXME: handle error
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

// Static serves static files.
func Static(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, staticBase+filepath.Base(r.URL.Path))
}

func method(r *http.Request) string {
	form := r.FormValue("method")
	if form != "" {
		return strings.ToUpper(form)
	}
	return r.Method
}
