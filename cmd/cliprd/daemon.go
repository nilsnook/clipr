package main

import (
	"context"
	"internal/data"
	"log"
	"os"
	"os/signal"
	"syscall"
)

type daemon struct {
	ctx     context.Context
	cancel  context.CancelFunc
	log     *log.Logger
	sigchan chan os.Signal
	event   *event
	db      *data.CliprDB
}

func newDaemon(log *log.Logger) *daemon {
	// init context
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	// init database
	db, err := data.NewCliprDB()
	if err != nil {
		log.Fatalln(err)
	}

	// init channel for system signals
	sigchan := make(chan os.Signal, 1)

	// init events
	e := &event{}
	e.copy = make(chan string, 1)

	// create and return daemon
	return &daemon{
		ctx:     ctx,
		cancel:  cancel,
		log:     log,
		sigchan: sigchan,
		event:   e,
		db:      db,
	}
}

func (d *daemon) handleSignals(sigs []os.Signal) {
	if len(sigs) > 0 {
		// setup notifications for daemon signals
		signal.Notify(d.sigchan, sigs...)

		for {
			select {
			case s := <-d.sigchan:
				switch s {
				case syscall.SIGINT, syscall.SIGTERM:
					d.log.Println("Got SIGINT/SIGTERM, exiting.")
					d.cancel()
					os.Exit(1)
				case syscall.SIGHUP:
					d.log.Println("Got SIGHUP, restarting with changed configuration.")
					// nothing to do here, because as af now there is no configuration
				}
			case <-d.ctx.Done():
				d.log.Fatalln("Done.")
			}
		}

	} else {
		d.log.Println("(0) signals to handle")
	}
}

// func (d *daemon) getLatestTextFromClipboard() string {
// 	// get latest copied text from system clipboard
// 	txt, err := exec.Command("xsel", "-ob").Output()
// 	if err != nil {
// 		d.log.Fatalln(err)
// 	}
// 	// rofiEncode - replace newline (\n) or carriage return (\r) with '\xA0'
// 	// before writing to database
// 	return data.RofiEncode(string(txt))
// }

func (d *daemon) run() {
	for {
		select {
		case <-d.ctx.Done():
			d.log.Fatalln("Done.")
		case txt := <-d.event.copy:
			// save text to clipr database
			err := d.db.Insert(txt)
			if err != nil {
				d.log.Fatalln(err)
			}
		}
	}
}

func (d *daemon) cleanup() {
	signal.Stop(d.sigchan)
	close(d.event.copy)
	d.cancel()
}
