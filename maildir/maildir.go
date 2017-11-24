// Package maildir implements parsing of maildir folders.
package maildir

import (
	"log"
	"net/mail"
	"os"
	"path/filepath"
	"time"
)

type Message struct {
	Date    time.Time
	From    mail.Address
	To      mail.Address
	Subject string
	Body    string
}

var Messages []*Message

func init() {
	mails, err := filepath.Glob("etc/mail/*:*")
	if err != nil {
		log.Panicf("error reading maildir: %s", err)
	}
	Messages = make([]*Message, len(mails))
	for i, m := range mails {
		file, err := os.Open(m)
		if err != nil {
			log.Panicf("error opening mail: %s", err)
		}
		msg, err := mail.ReadMessage(file)
		if err != nil {
			log.Panicf("error parsing mail: %s", err)
		}
		Messages[i] = newMessage(msg)
	}
}

func newMessage(msg *mail.Message) *Message {
	date, err := msg.Header.Date()
	from, err := msg.Header.AddressList("From")
	to, err := msg.Header.AddressList("To")
	if err != nil {
		log.Panicf("error parsing header: %s", err)
	}

	return &Message{date, *from[0], *to[0], msg.Header.Get("Subject"), msg.Header.Get("Body")}
}
