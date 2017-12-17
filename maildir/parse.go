package maildir

import (
	"io/ioutil"
	"log"
	"mime"
	"mime/multipart"
	"mime/quotedprintable"
	"net/mail"
	"strings"
)

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

			raw, err := ioutil.ReadAll(p)
			contentType, _, err := mime.ParseMediaType(p.Header.Get("Content-Type"))
			if err != nil {
				log.Println("[ERROR] maildir:", err)
			}

			if contentType == "text/html" {
				body.HTML = string(raw)
			} else {
				body.Plain = string(raw)
			}
		}
	} else {
		var raw []byte
		if strings.ToLower(msg.Header.Get("Content-Transfer-Encoding")) == "quoted-printable" {
			raw, err = ioutil.ReadAll(quotedprintable.NewReader(msg.Body))
		} else {
			raw, err = ioutil.ReadAll(msg.Body)
		}

		contentType, _, err := mime.ParseMediaType(msg.Header.Get("Content-Type"))
		if err != nil {
			log.Println("[ERROR] maildir:", err)
		}

		if contentType == "text/html" {
			body.HTML = string(raw)
		} else {
			body.Plain = string(raw)
		}
	}

	return
}
