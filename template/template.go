// Package template loads, parses, and renders HTML templates.
package template

import (
	"html/template"
	"io"
	"log"
)

const base = "template/base.tmpl"

var templates = make(map[string]*template.Template)

func init() {
	templates["index.html"] = template.Must(template.New("base").ParseFiles(base, "template/index.tmpl"))
	templates["compose.html"] = template.Must(template.New("base").ParseFiles(base, "template/compose.tmpl"))
	templates["message.html"] = template.Must(template.New("base").ParseFiles(base, "template/message.tmpl"))
}

// Render renders a template with optional data and writes the result to w.
func Render(w io.Writer, tmpl string, data interface{}) {
	template, ok := templates[tmpl]
	if !ok {
		log.Panicf("no template %s", tmpl)
	}
	err := template.Execute(w, data)
	if err != nil {
		log.Panicf("error rendering %s: %s", tmpl, err)
	}
}
