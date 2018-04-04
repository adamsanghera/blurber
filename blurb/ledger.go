package blurb

import (
	"log"
	"sync"
	"time"
)

// Ledger is the interface for a service that maintains
// a record of blurbs (social media posts)
type Ledger interface {
	AddBlurb(creator string, b Blurb)
	RemoveBlurb(creator string, b Blurb)
}

type FeedCache struct {
	timestamp  time.Time
	sortedFeed []Blurb
}

// LocalLedger is an implementation of Ledger, making use
// of in-memory data structures to record blurbs
type LocalLedger struct {
	// Map from UID -> map [BID] -> Blurb
	bidCounter int
	bidMutex   sync.Mutex

	ledger    map[int]*BlurbBox
	feedCache map[int]FeedCache
}

// NewLocalLedger returns a new and initialized LocalLedger
func NewLocalLedger() *LocalLedger {
	log.Printf("BLURB-LEDGER: Initializing")
	return &LocalLedger{
		ledger:     make(map[int]*BlurbBox),
		feedCache:  make(map[int]FeedCache),
		bidCounter: 0,
		bidMutex:   sync.Mutex{},
	}
}

func (ll *LocalLedger) measureBID() int {
	ll.bidMutex.Lock()

	ret := ll.bidCounter
	ll.bidCounter++

	ll.bidMutex.Unlock()

	return ret
}