package data

import "time"

type Val string

type Meta struct {
	LastModified time.Time `json:"last_modified"`
}

type Entry struct {
	Val  Val  `json:"val"`
	Meta Meta `json:"meta"`
}

type Clipboard struct {
	List []Entry `json:"list"`
}
