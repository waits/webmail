package handler

import (
	"context"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"gopkg.in/hlandau/passlib.v1"
)

type user struct {
	pass  string
	tries int
}

type key int

const userKey key = 0

var users map[string]*user

// LoadPasswd reads the passwd file at path into memory.
func LoadPasswd(path string) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	lines := strings.Split(string(content), "\n")
	users = make(map[string]*user, len(lines))
	for _, l := range lines {
		secs := strings.Split(strings.TrimSpace(l), ":")
		if len(secs) < 2 {
			continue
		}
		users[secs[0]] = &user{pass: strings.TrimPrefix(secs[1], "{SHA256-CRYPT}")}
	}
}

// WithAuth wraps an http.HandleFunc to require authentication.
func WithAuth(fn func(w http.ResponseWriter, r *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if un, pass, ok := r.BasicAuth(); ok {
			u, ok := users[un]
			if ok && u.tries < 10 {
				if _, err := passlib.Verify(pass, u.pass); err == nil {
					log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL.Path)
					ctx := context.WithValue(r.Context(), userKey, u)
					fn(w, r.WithContext(ctx))
					u.tries = 0
					return
				}
				u.tries++
				log.Printf("%s tries %d", un, u.tries)
			}
		}

		log.Printf("%s %s %s %d", r.RemoteAddr, r.Method, r.URL.Path, http.StatusUnauthorized)
		w.Header().Set("WWW-Authenticate", "Basic realm=\"webmail\"")
		http.Error(w, "401 unauthorized", http.StatusUnauthorized)
	}
}
