package main

import (
	"os/exec"
)

type event struct {
	copy chan string
}

func (e *event) subscribeToCopy(handler func() string) {
	for {
		clipnotifyCmd := exec.Command("clipnotify")
		_, err := clipnotifyCmd.Output()
		if err == nil {
			e.copy <- handler()
		}
	}
}
