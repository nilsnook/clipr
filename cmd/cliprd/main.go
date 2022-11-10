package main

import (
	"log"
	"os"
	"syscall"
)

func main() {
	// set log output to STDOUT explicitly
	log := log.New(os.Stdout, "CLIPR_D\t", log.LstdFlags|log.Lshortfile)

	// create new daemon
	d := newDaemon(log)
	defer d.cleanup()

	// handle daemon signals
	sigs := []os.Signal{
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGHUP,
	}
	go d.handleSignals(sigs)
	// subscribe to copy event
	go d.event.subscribeToCopy(d.getLatestTextFromClipboard)

	// run daemon
	d.run()
}
