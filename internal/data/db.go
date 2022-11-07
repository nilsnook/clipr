package data

import (
	"errors"
	"os"
	"path"
	"time"
)

const (
	DB_DIR  = ".local/share/clipr"
	DB_NAME = "clipr.db"
)

type CliprDB struct {
	Dir       string
	Name      string
	Clipboard Clipboard
}

func NewCliprDB(dir string) (*CliprDB, error) {
	// get user home directory
	userhomedir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	// resolve db dir with user home dir
	dir = resolveDir(userhomedir, dir)
	// retuen new instance
	return &CliprDB{
		Dir:       dir,
		Name:      DB_NAME,
		Clipboard: Clipboard{},
	}, nil
}

func (db *CliprDB) getDBFile(flag int) (f *os.File, err error) {
	DB := path.Join(db.Dir, db.Name)
	f, err = os.OpenFile(DB, flag, 0644)
	return
}

func (db *CliprDB) CreateClipboardIfNotExists() (err error) {
	var f *os.File
	defer f.Close()

	DB := path.Join(db.Dir, db.Name)
	if _, err = os.Stat(DB); errors.Is(err, os.ErrNotExist) {
		// create dir
		err = createDirIfNotExists(db.Dir)
		if err != nil {
			return
		}

		// create db file
		f, err = os.OpenFile(DB, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return
		}
		defer f.Close()

		// create a clipboard and save to file
		x := Clipboard{
			List: []Entry{},
		}
		err = writeJSON(f, x)
		if err != nil {
			return
		}
	}
	return
}

func (db *CliprDB) Read() (err error) {
	// get db file
	f, err := db.getDBFile(os.O_RDONLY)
	if err != nil {
		return
	}
	defer f.Close()

	// read from db file
	db.Clipboard, err = readJSON(f)
	if err != nil {
		return
	}
	return
}

func (db *CliprDB) Write(txt string) (err error) {
	// prepare entry from text
	e := Entry{
		Val: Val(txt),
		Meta: Meta{
			LastModified: time.Now(),
		},
	}
	// append new entry to existing clipboard
	db.Clipboard.List = append(db.Clipboard.List, e)
	// create new set from clipboard entries
	s := NewSet(db.Clipboard.List...)
	// get unique set of entries
	db.Clipboard.List = s.Entries()

	// get db file
	f, err := db.getDBFile(os.O_WRONLY)
	if err != nil {
		return
	}
	defer f.Close()

	// write changes to db file
	err = writeJSON(f, db.Clipboard)
	if err != nil {
		return
	}
	return
}

func (db *CliprDB) Delete(txt string) (err error) {
	// prepare entry from text
	e := Entry{
		Val: Val(txt),
	}
	// create a set of existing clipboard entries
	s := NewSet(db.Clipboard.List...)
	// delete the required entry
	s.Delete(e)
	// get new list of entries from set
	db.Clipboard.List = s.Entries()

	// get db file
	f, err := db.getDBFile(os.O_WRONLY)
	if err != nil {
		return
	}
	defer f.Close()

	// write changes to db file
	err = writeJSON(f, db.Clipboard)
	if err != nil {
		return
	}
	return
}
