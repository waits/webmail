// Package maildir implements parsing of maildir folders.
package maildir

import (
	"crypto/sha256"
	"encoding/hex"
	"log"
	"net/mail"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/fsnotify/fsnotify"
)

// Message represents a single mail message.
type Message struct {
	Date     time.Time
	ID       string
	From     []*mail.Address
	FromName string
	To       []*mail.Address
	Subject  string
	Body     Body
	Flag     Flag
	path     string
}

// Time returns a formatted time string.
func (m *Message) Time() string {
	return m.Date.Local().Format("Jan 2 3:04pm")
}

// Body holds multipart mail bodies for a single message.
type Body struct {
	Plain string
	HTML  string
}

// Messages is a map of IDs/names to messages.
var Messages map[string]*Message
var fileMap map[string]*Message

// ByDate implements sort.Interface for []Message based on the Date field.
type ByDate []*Message

func (a ByDate) Len() int           { return len(a) }
func (a ByDate) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByDate) Less(i, j int) bool { return a[i].Date.Before(a[j].Date) }

// Watch watches dir for changes and populates Messages based on received events.
func Watch(dir string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	// Watcher event loop
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

	dir = filepath.Join(dir, "cur")

	// Initialize Messages with dir contents
	msgs, err := filepath.Glob(dir + "/*:*")
	if err != nil {
		log.Fatalln("[ERROR] maildir:", err)
	}
	Messages = make(map[string]*Message, len(msgs))
	fileMap = make(map[string]*Message, len(msgs))
	for _, m := range msgs {
		openMessage(m)
	}

	// Start watching dir
	// FIXME: race condition?
	err = watcher.Add(dir)
	if err != nil {
		log.Fatal(err)
	}
}

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

// Sorted returns Messages sorted by date descending.
func Sorted() []*Message {
	var messages ByDate
	for _, m := range Messages {
		messages = append(messages, m)
	}
	sort.Sort(sort.Reverse(messages))
	return messages
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
	body, err := parseBody(msg)
	flag, err := parseFlags(name)
	if err != nil {
		log.Println("[ERROR] maildir:", err)
	}

	fromName := "Unknown Sender"
	if len(from) > 0 {
		fromName = from[0].Name
	}

	checksum := sha256.Sum256([]byte(msg.Header.Get("Message-ID")))
	id := hex.EncodeToString(checksum[:8])

	return &Message{date, id, from, fromName, to, msg.Header.Get("Subject"), body, flag, name}
}
