package blurb

import (
	"time"
)

type Blurb struct {
	Content     string
	Timestamp   string
	Time        time.Time
	BID         int // immutable
	CreatorName string
}
