package store

import "time"

type Item struct {
	Value      any
	Expiration int64
}

func (i Item) Expired() bool {
	if i.Expiration == 0 {
		return false
	}
	return time.Now().UnixNano() > i.Expiration
}
