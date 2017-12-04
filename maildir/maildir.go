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

const dir = "tmp/inbox/cur"

// Message represents a single mail message.
type Message struct {
	Date     time.Time
	ID       string
	From     []*mail.Address
	FromName string
	To       []*mail.Address
	Subject  string
	Body     string
	path     string
}

// Messages is a map of IDs/names to messages.
var Messages map[string]*Message
var fileMap map[string]*Message

// DeleteMessage deletes a message from Messages and the filesystem.
func DeleteMessage(key string) error {
	msg, ok := Messages[key]
	if ok {
		delete(Messages, msg.ID)
		delete(fileMap, msg.path)
		return os.Remove(msg.path)
	}
	return nil
}

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
					DeleteMessage(event.Name)
				}
			case err := <-watcher.Errors:
				log.Println("[ERROR] maildir: watcher error:", err)
			}
		}
	}()

	initMessages() // FIXME: race condition?

	err = watcher.Add(dir)
	if err != nil {
		log.Fatal(err)
	}
}

func initMessages() {
	mails, err := filepath.Glob(dir + "/*:*")
	if err != nil {
		log.Fatalln("[ERROR] maildir:", err)
	}
	Messages = make(map[string]*Message, len(mails))
	fileMap = make(map[string]*Message, len(mails))
	for _, m := range mails {
		openMessage(m)
	}
}

func openMessage(name string) {
	file, err := os.Open(name)
	if err != nil {
		log.Println("[ERROR] maildir:", err)
		return
	}
	msg, err := mail.ReadMessage(file)
	if err != nil {
		log.Println("[ERROR] maildir:", err)
		return
	}
	message := newMessage(msg, name)
	fileMap[message.path] = message
	Messages[message.ID] = message
}

func newMessage(msg *mail.Message, name string) *Message {
	date, err := msg.Header.Date()
	from, err := msg.Header.AddressList("From")
	to, err := msg.Header.AddressList("To")
	body, err := ioutil.ReadAll(msg.Body)
	if err != nil {
		log.Fatalln("[ERROR] maildir:", err)
	}

	fromName := "Unknown Sender"
	if len(from) > 0 {
		fromName = from[0].Name
	}

	checksum := sha256.Sum256([]byte(msg.Header.Get("Message-ID")))
	id := hex.EncodeToString(checksum[:8])

	return &Message{date, id, from, fromName, to, msg.Header.Get("Subject"), string(body), name}
}
