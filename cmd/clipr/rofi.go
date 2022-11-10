package main

import (
	"encoding/json"
	"errors"
	"internal/data"
	"os"
	"strconv"
)

const (
	// Enter
	SELECT = 1
	// kb-custom-1
	DELETE = 10
)

type info struct {
	Id int `json:"id"`
}

type state struct {
	val  int
	info info
	arg  string
}

type rofi struct {
	state state
}

func newRofi() *rofi {
	r := &rofi{}
	r.initState()
	return r
}

func (r *rofi) initState() {
	if val := os.Getenv("ROFI_RETV"); val != "" {
		// c.infolog.Printf("State: %s", val)
		r.state.val, _ = strconv.Atoi(val)
	}

	if info := os.Getenv("ROFI_INFO"); info != "" {
		// c.infolog.Printf("Info: %s", info)
		json.Unmarshal([]byte(info), &r.state.info)
	}

	args := os.Args
	if len(args) > 1 {
		for k, v := range args {
			if k == 1 {
				// c.infolog.Printf("Arg: %s", v)
				r.state.arg = v
			}
		}
	}
}

func (r *rofi) getSelection() (string, error) {
	txt := data.RofiDecode(r.state.arg)
	if len(txt) == 0 {
		return txt, errors.New("Selection empty! Failed to copy to clipboard.")
	}
	return txt, nil
}
