package main

import (
	"fmt"
	"log"
	"os"
	"path"
)

func main() {
	LOGFILE := path.Join(os.TempDir(), "clipr.log")
	f, err := os.OpenFile(LOGFILE, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()

	// create new clipr with specified log file
	c := newClipr(f)

	// handle rofi events
	switch c.rofi.state.val {
	case SELECT:
		c.copySelection()
	case DELETE:
		c.deleteSelection()
	}

	// use hot keys
	// for events like delete
	fmt.Println("\000use-hot-keys\x1ftrue")
	// render clipboard
	c.renderClipboard()
}
