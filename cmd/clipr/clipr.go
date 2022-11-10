package main

import (
	"fmt"
	"internal/data"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
)

type clipboard struct {
	List []data.Entry `json:"list"`
}

type clipr struct {
	infolog  *log.Logger
	errorlog *log.Logger
	db       *data.CliprDB
	rofi     *rofi
}

func newClipr(f *os.File) *clipr {
	infolog := log.New(f, "INFO\t", log.LstdFlags)
	errorlog := log.New(f, "ERROR\t", log.LstdFlags|log.Lshortfile)
	db, err := data.NewCliprDB("~/.local/share/clipr")
	if err != nil {
		log.Fatal(err)
	}
	rofi := newRofi()
	return &clipr{
		infolog:  infolog,
		errorlog: errorlog,
		db:       db,
		rofi:     rofi,
	}
}

func (c *clipr) initClipboard() {
	var err error
	// create one if does not exists
	err = c.db.CreateClipboardIfNotExists()
	if err != nil {
		c.errorlog.Fatalln(err)
	}
	// read from clipboard
	err = c.db.Read()
	if err != nil {
		c.errorlog.Fatalln(err)
	}
}

func (c *clipr) getLatestTextFromClipboard() {
	// get latest copied text from system clipboard
	t, err := exec.Command("xsel", "-ob").Output()
	if err != nil {
		c.errorlog.Fatalln(err)
	}
	// replace newline (\n) or carriage return (\r) with '\xA0'
	// before writing entry into database
	txt := data.RofiEncode(string(t))
	c.db.Write(txt)
}

func (c *clipr) renderClipboard() {
	for _, e := range c.db.Clipboard.List {
		fmt.Println(e.Val)
	}
}

func (c *clipr) copySelection() {
	// get selection from rofi
	txt, err := c.rofi.getSelection()
	if err != nil {
		c.errorlog.Fatalln("Selection empty! Failed to copy to clipboard.")
		return
	}

	// copy selection to system clipboard
	pr, pw := io.Pipe()
	xselcmd := exec.Command("xsel", "-ib")
	xselcmd.Stdin = pr
	go func() {
		defer pw.Close()
		io.Copy(pw, strings.NewReader(txt))
	}()
	err = xselcmd.Run()
	if err != nil {
		c.errorlog.Fatalln(err)
	}
}

func (c *clipr) deleteSelection() {
	// get selection from rofi
	txt, err := c.rofi.getSelection()
	if err != nil {
		c.errorlog.Fatalln("Selection empty! Failed to delete clipboard entry.")
		return
	}

	// delete selection from db
	c.db.Delete(txt)

	// if the last entry is deleted
	// clear system clipboard as well
	if len(c.db.Clipboard.List) == 0 {
		err := exec.Command("xsel", "-cb").Run()
		if err != nil {
			c.errorlog.Fatalln(err)
		}
	}
}
