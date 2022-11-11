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
	db, err := data.NewCliprDB()
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

func (c *clipr) renderClipboard() {
	// read clipboard entries from db
	clipboard, err := c.db.Read()
	if err != nil {
		c.errorlog.Fatalln(err)
	}
	// traverse through each entry and print
	for _, e := range clipboard.List {
		fmt.Println(e.Val)
	}
}

func (c *clipr) copySelection() {
	// get selection from rofi
	sel, err := c.rofi.getSelection()
	if err != nil {
		c.errorlog.Fatalln("Selection empty! Failed to copy to clipboard.")
		return
	}
	// rofi decode the selection
	txt := data.RofiDecode(sel)

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
	sel, err := c.rofi.getSelection()
	if err != nil {
		c.errorlog.Fatalln("Selection empty! Failed to delete clipboard entry.")
		return
	}

	// delete selection from db
	cbsize, _ := c.db.Delete(sel)

	// if the last entry is deleted
	// clear system clipboard as well
	if cbsize == 0 {
		err := exec.Command("xsel", "-cb").Run()
		if err != nil {
			c.errorlog.Fatalln(err)
		}
	}
}
