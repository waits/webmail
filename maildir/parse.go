package maildir

import (
	"encoding/base64"
	"io"
	"io/ioutil"
	"log"
	"mime"
	"mime/multipart"
	"mime/quotedprintable"
	"net/mail"
	"strings"
)

type header interface {
	Get(string) string
}

// Parses multipart and quoted-printable mail bodies.
func parseBody(msg *mail.Message) (body Body, err error) {
	mediaType, params, err := mime.ParseMediaType(msg.Header.Get("Content-Type"))
	if err != nil {
		log.Println("[ERROR] maildir:", err)
	}

	if strings.HasPrefix(mediaType, "multipart/") {
		mr := multipart.NewReader(msg.Body, params["boundary"])
		for {
			p, err := mr.NextPart()
			if err != nil {
				break
			}

			decodePart(p.Header, p, &body)
		}
	} else {
		decodePart(msg.Header, msg.Body, &body)
	}

	return
}

func decodePart(h header, body io.Reader, out *Body) {
	var raw []byte
	var err error

	if strings.ToLower(h.Get("Content-Transfer-Encoding")) == "quoted-printable" {
		raw, err = ioutil.ReadAll(quotedprintable.NewReader(body))
	} else if strings.ToLower(h.Get("Content-Transfer-Encoding")) == "base64" {
		raw, err = ioutil.ReadAll(base64.NewDecoder(base64.StdEncoding, body))
	} else {
		raw, err = ioutil.ReadAll(body)
	}

	contentType, _, err := mime.ParseMediaType(h.Get("Content-Type"))
	if err != nil {
		log.Println("[ERROR] maildir:", err)
	}

	if contentType == "text/html" {
		out.HTML = string(raw)
	} else {
		out.Plain = string(raw)
	}
}
