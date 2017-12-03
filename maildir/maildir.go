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

	"github.com/fsnotify/fsnotify"
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

// fileMap is a map of filenames to message IDs.
var fileMap map[string]string

func init() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for {
			select {
			case event := <-watcher.Events:
				switch event.Op {
				case fsnotify.Create, fsnotify.Write:
					log.Printf("maildir: opening %s", event.Name)
					openMessage(event.Name)
				case fsnotify.Remove, fsnotify.Rename:
					log.Printf("maildir: removing %s", event.Name)
					removeMessage(event.Name)
				}
			case err := <-watcher.Errors:
				log.Println("[ERROR] maildir: watcher error:", err)
			}
		}
	}()

	initMessages() // FIXME: race condition?

	err = watcher.Add("etc/mail")
	if err != nil {
		log.Fatal(err)
	}
}

func initMessages() {
	mails, err := filepath.Glob("etc/mail/*:*")
	if err != nil {
		log.Fatalln("[ERROR] maildir:", err)
	}
	Messages = make(map[string]*Message, len(mails))
	fileMap = make(map[string]string, len(mails))
	for _, m := range mails {
		openMessage(m)
	}
}

func openMessage(m string) {
	file, err := os.Open(m)
	if err != nil {
		log.Println("[ERROR] maildir:", err)
		return
	}
	msg, err := mail.ReadMessage(file)
	if err != nil {
		log.Println("[ERROR] maildir:", err)
		return
	}
	message := newMessage(msg)
	Messages[message.ID] = message
	fileMap[m] = message.ID // TODO: split name at colon.
}

func removeMessage(m string) {
	id, ok := fileMap[m]
	if ok {
		delete(fileMap, m)
		delete(Messages, id)
	}
}

func newMessage(msg *mail.Message) *Message {
	date, err := msg.Header.Date()
	from, err := msg.Header.AddressList("From")
	to, err := msg.Header.AddressList("To")
	body, err := ioutil.ReadAll(msg.Body)
	if err != nil {
		log.Fatalln("[ERROR] maildir:", err)
	}

	checksum := sha256.Sum256([]byte(msg.Header.Get("Message-ID")))
	id := hex.EncodeToString(checksum[:8])

	return &Message{date, id, *from[0], *to[0], msg.Header.Get("Subject"), string(body)}
}
