// Package maildir implements parsing of maildir folders.
package maildir

import (
	"crypto/sha256"
	"encoding/hex"
	"io/ioutil"
	"log"
	"net/mail"
	"os"
	"path/filepath"
	"time"
)

// Message represents a single mail message.
type Message struct {
	Date    time.Time
	ID      string
	From    mail.Address
	To      mail.Address
	Subject string
	Body    string
}

// Messages is a map of IDs to messages.
var Messages map[string]*Message

func init() {
	mails, err := filepath.Glob("etc/mail/*:*")
	if err != nil {
		log.Panicf("error reading maildir: %s", err)
	}
	Messages = make(map[string]*Message, len(mails))
	for _, m := range mails {
		file, err := os.Open(m)
		if err != nil {
			log.Panicf("error opening mail: %s", err)
		}
		msg, err := mail.ReadMessage(file)
		if err != nil {
			log.Panicf("error parsing mail: %s", err)
		}
		message := newMessage(msg)
		Messages[message.ID] = message
	}
}

func newMessage(msg *mail.Message) *Message {
	date, err := msg.Header.Date()
	from, err := msg.Header.AddressList("From")
	to, err := msg.Header.AddressList("To")
	body, err := ioutil.ReadAll(msg.Body)
	if err != nil {
		log.Panicf("error parsing header: %s", err)
	}

	checksum := sha256.Sum256([]byte(msg.Header.Get("Message-ID")))
	id := hex.EncodeToString(checksum[:8])

	return &Message{date, id, *from[0], *to[0], msg.Header.Get("Subject"), string(body)}
}
