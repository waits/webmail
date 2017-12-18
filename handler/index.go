// Package handler implements HTTP request handlers.
package handler

import (
	"net/http"
	"net/smtp"
	"path/filepath"
	"strings"

	"github.com/waits/webmail/maildir"
	"github.com/waits/webmail/template"
)

const staticBase = "static/"

// Index serves the home page.
func Index(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	template.Render(w, "index.html", maildir.Sorted())
}

// Compose serves the compose mail page.
func Compose(w http.ResponseWriter, r *http.Request) {
	template.Render(w, "compose.html", nil)
}

// Send sends a mail message.
func Send(w http.ResponseWriter, r *http.Request) {
	auth := smtp.PlainAuth("", "", "", "localhost")
	from := r.FormValue("from")
	to := strings.Split(r.FormValue("to"), ", ")
	body := []byte(r.FormValue("body"))
	err := smtp.SendMail("localhost:25", auth, from, to, body)
	if err != nil {
		http.Error(w, "failed to connect to SMTP server", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Message serves the message detail page.
func Message(w http.ResponseWriter, r *http.Request) {
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
		msg.Flag.Write('S')
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
