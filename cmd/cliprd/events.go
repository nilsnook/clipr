package main

import (
	"context"
	"internal/data"

	"golang.design/x/clipboard"
)

type event struct {
	copy chan string
}

// func (e *event) subscribeToCopy(handler func() string) {
// 	for {
// 		clipnotifyCmd := exec.Command("clipnotify")
// 		_, err := clipnotifyCmd.Output()
// 		if err == nil {
// 			e.copy <- handler()
// 		}
// 	}
// }

// func (e *event) subscribeToCopy(handler func() string) {
// 	clipnotifyCmd := exec.Command("clipnotify")
// 	stdout, err := clipnotifyCmd.StdoutPipe()
// 	if err == nil {
// 		scanner := bufio.NewScanner(stdout)
// 		err = clipnotifyCmd.Start()
// 		if err == nil {
// 			for scanner.Scan() {
// 				e.copy <- handler()
// 			}
// 		}
// 		if scanner.Err() != nil {
// 			clipnotifyCmd.Process.Kill()
// 			clipnotifyCmd.Wait()
// 		}
// 	}
// 	clipnotifyCmd.Wait()
// }

func (e *event) subscribeToCopy() {
	ch := clipboard.Watch(context.TODO(), clipboard.FmtText)
	for txt := range ch {
		// passing down the latest copied text to the event's copy channel.
		// rofiEncode - replace newline (\n) or carriage return (\r) with '\xA0'
		// before writing to database.
		e.copy <- data.RofiEncode(string(txt))
	}
}
