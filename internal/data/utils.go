package data

import (
	"encoding/json"
	"errors"
	"os"
	"path"
	"strings"
	"unicode/utf8"
)

func resolveDir(homeDir, dir string) string {
	if strings.HasPrefix(dir, "~/") {
		dir = strings.TrimLeft(dir, "~/")
		dir = path.Join(homeDir, dir)
	} else if strings.HasPrefix(dir, "$HOME/") {
		dir = strings.TrimLeft(dir, "$HOME/")
		dir = path.Join(homeDir, dir)
	}
	return dir
}

func createDirIfNotExists(path string) (err error) {
	if _, err = os.Stat(path); errors.Is(err, os.ErrNotExist) {
		err = os.MkdirAll(path, 0755)
	}
	return
}

func writeJSON(f *os.File, data Clipboard) (err error) {
	err = json.NewEncoder(f).Encode(data)
	return
}

func readJSON(f *os.File) (data Clipboard, err error) {
	err = json.NewDecoder(f).Decode(&data)
	return
}

func RofiEncode(txt string) (enctxt string) {
	rtxt := make([]rune, 0, utf8.RuneCountInString(txt))
	for _, ch := range txt {
		if ch == '\n' || ch == '\r' {
			rtxt = append(rtxt, '\xA0')
		} else {
			rtxt = append(rtxt, ch)
		}
	}
	enctxt = string(rtxt)
	return
}

func RofiDecode(enctxt string) (txt string) {
	ctxt := make([]rune, 0, utf8.RuneCountInString(enctxt))
	for _, ch := range enctxt {
		if ch == '\xA0' {
			ctxt = append(ctxt, '\n')
		} else {
			ctxt = append(ctxt, ch)
		}
	}
	txt = string(ctxt)
	return
}
