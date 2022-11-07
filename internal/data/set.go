package data

import "sort"

type set map[Val]Meta

func NewSet(entries ...Entry) (s set) {
	s = set{}
	for _, e := range entries {
		s[e.Val] = e.Meta
	}
	return
}

func (s *set) Add(entries ...Entry) {
	for _, e := range entries {
		(*s)[e.Val] = e.Meta
	}
}

func (s *set) Delete(e Entry) {
	delete(*s, e.Val)
}

func (s set) Entries() []Entry {
	entries := make([]Entry, 0, len(s))
	for k, v := range s {
		e := Entry{
			Val:  k,
			Meta: v,
		}
		entries = append(entries, e)
	}
	// sort entries w.r.t last modified time
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Meta.LastModified.After(entries[j].Meta.LastModified)
	})
	return entries
}
