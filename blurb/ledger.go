package blurb

import (
	"log"
	"sync"
	"time"

	"github.com/adamsanghera/blurber/protobufs/dist/blurb"
)

// Ledger is the interface for a service that maintains
// a record of blurbs (social media posts)
type Ledger interface {
	AddBlurb(creator string, b blurb.Blurb)
	RemoveBlurb(creator string, b blurb.Blurb)
}

type FeedCache struct {
	timestamp  time.Time
	sortedFeed []blurb.Blurb
}

// LocalLedger is an implementation of Ledger, making use
// of in-memory data structures to record blurbs
type LocalLedger struct {
	// Map from UID -> map [BID] -> Blurb
	bidCounter int32
	bidMutex   sync.Mutex

	ledger    sync.Map // Stores BlurbBox objects
	feedCache sync.Map // Stores FeedCache objects
}

// NewLocalLedger returns a new and initialized LocalLedger
func NewLocalLedger() *LocalLedger {
	log.Printf("BLURB-LEDGER: Initializing")
	return &LocalLedger{
		ledger:     sync.Map{},
		feedCache:  sync.Map{},
		bidCounter: 0,
		bidMutex:   sync.Mutex{},
	}
}

func (ll *LocalLedger) measureBID() int32 {
	ll.bidMutex.Lock()

	ret := ll.bidCounter
	ll.bidCounter++

	ll.bidMutex.Unlock()

	return ret
}
