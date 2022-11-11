package data

import (
	"errors"
	"os"
	"path"
	"strings"
	"time"
)

const (
	DB_DIR  = "~/.local/share/clipr"
	DB_NAME = "clipr.db"
)

type CliprDB struct {
	Dir  string
	Name string
}

func NewCliprDB(dirs ...string) (*CliprDB, error) {
	// get user home directory
	userhomedir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	// check for valid dir value(s)
	dir := DB_DIR
	if len(dirs) > 0 {
		if len(strings.Trim(dirs[0], " ")) > 0 {
			dir = dirs[0]
		}
	}

	// resolve db dir for user home dir, if user home path exists
	dir = resolveDir(userhomedir, dir)

	// create db dir if not exists
	if _, err = os.Stat(dir); errors.Is(err, os.ErrNotExist) {
		err = createDirIfNotExists(dir)
		if err != nil {
			return nil, err
		}
	}

	// retuen new instance
	return &CliprDB{
		Dir:  dir,
		Name: DB_NAME,
	}, nil
}

func (db *CliprDB) Read() (cb Clipboard, err error) {
	// get db file
	DB := path.Join(db.Dir, db.Name)
	f, err := os.OpenFile(DB, os.O_RDONLY, 0644)
	if err != nil {
		return
	}
	defer f.Close()

	// read from db file
	cb, err = readJSON(f)
	if err != nil {
		return
	}
	return
}

func (db *CliprDB) Write(clipboard Clipboard) (err error) {
	// get db file
	DB := path.Join(db.Dir, db.Name)
	f, err := os.OpenFile(DB, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return
	}
	defer f.Close()

	// write changes to db file
	err = writeJSON(f, clipboard)
	if err != nil {
		return
	}
	return
}

func (db *CliprDB) Insert(txt string) (err error) {
	// return if empty text
	if len(txt) == 0 {
		return
	}

	// prepare entry from text
	e := Entry{
		Val: Val(txt),
		Meta: Meta{
			LastModified: time.Now(),
		},
	}
	// read from db before every write to maintain sync
	// with database and in-memory struct
	clipboard, _ := db.Read()
	// append new entry to existing clipboard
	clipboard.List = append(clipboard.List, e)
	// create new set from clipboard entries
	s := NewSet(clipboard.List...)
	// get unique set of entries
	clipboard.List = s.Entries()

	// write clipboard to DB
	err = db.Write(clipboard)
	return
}

func (db *CliprDB) Delete(txt string) (size int, err error) {
	// prepare entry from text
	e := Entry{
		Val: Val(txt),
	}
	// read from db before every write to maintain sync
	// with database and in-memory struct
	clipboard, _ := db.Read()
	// create a set of existing clipboard entries
	s := NewSet(clipboard.List...)
	// delete the required entry
	s.Delete(e)
	// get new list of entries from set
	clipboard.List = s.Entries()
	// set current clipboard size
	size = len(clipboard.List)

	// get db file
	err = db.Write(clipboard)
	return
}
