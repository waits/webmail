package maildir

import (
	"errors"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

const sep = ":2,"

// Flag holds flags parsed from maildir filenames.
type Flag struct {
	flags map[rune]bool // FIXME: use sync.Map for thread safety.
	path  string
}

// Get gets the boolean value of flag r.
func (f *Flag) Get(r rune) bool {
	return f.flags[r]
}

// Not is the inverse of Get. Useful for templates.
func (f *Flag) Not(r rune) bool {
	return !f.flags[r]
}

// Write writes flag r to the filesystem if it doesn't already exist.
func (f *Flag) Write(r rune) error {
	dir, file := filepath.Split(f.path)
	flstr := strings.Split(file, sep)
	if len(flstr) != 2 {
		return errors.New("invalid maildir filename")
	}

	chars := strings.Split(flstr[1], "")
	if strings.ContainsRune(flstr[1], r) {
		return nil
	}

	chars = append(chars, string(r))
	sort.Strings(chars)
	out := filepath.Join(dir, flstr[0]+sep+strings.Join(chars, ""))

	return os.Rename(f.path, out)
}

// Parses char flags from a maildir filename.
func parseFlags(path string) (Flag, error) {
	_, file := filepath.Split(path)
	flstr := strings.Split(file, sep)
	if len(flstr) != 2 {
		return Flag{}, errors.New("invalid maildir filename")
	}

	flmap := make(map[rune]bool, 6)
	for _, ch := range flstr[1] {
		flmap[ch] = true
	}

	return Flag{flmap, path}, nil
}
