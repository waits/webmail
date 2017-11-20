// Package template
package template

import (
	"html/template"
	"io"
	"log"
)

const base = "template/base.tmpl"

var templates map[string]*template.Template = make(map[string]*template.Template)

func init() {
	templates["index.html"] = template.Must(template.New("base").ParseFiles(base, "template/index.tmpl"))
	templates["compose.html"] = template.Must(template.New("base").ParseFiles(base, "template/compose.tmpl"))
}

func Render(w io.Writer, tmpl string, data interface{}) {
	err := templates[tmpl].Execute(w, data)
	if err != nil {
		log.Panicf("error rendering %s: %s", tmpl, err)
	}
}
