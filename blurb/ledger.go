package blurb

import (
	"log"
	"sync"

	"github.com/adamsanghera/blurber-protobufs/dist/blurb"
)

// Ledger is the interface for a service that maintains
// a record of blurbs (social media posts), indexed by UIDs
type Ledger interface {
	AddBlurb(creator string, b blurb.Blurb)
	RemoveBlurb(creator string, b blurb.Blurb)
}

// LocalLedger is an implementation of Ledger, making use
// of in-memory data structures to record blurbs
type LocalLedger struct {
	// Map from UID -> map [BID] -> Blurb
	ledger     sync.Map // Stores BlurbBox objects
	bidCounter int32
	bidMutex   sync.Mutex

	feeds struct {
		cache sync.Map // Stores FeedCache objects

		// Config of feed
		length        int
		blurbsPerUser int
	}
}

// NewLocalLedger returns a new and initialized LocalLedger
func NewLocalLedger() *LocalLedger {
	log.Printf("BLURB-LEDGER: Initializing")
	return &LocalLedger{
		ledger:     sync.Map{},
		bidCounter: 0,
		bidMutex:   sync.Mutex{},

		feeds: struct {
			cache         sync.Map
			length        int
			blurbsPerUser int
		}{
			cache:         sync.Map{},
			length:        100,
			blurbsPerUser: 10,
		},
	}
}

func (ll *LocalLedger) measureBID() int32 {
	ll.bidMutex.Lock()

	ret := ll.bidCounter
	ll.bidCounter++

	ll.bidMutex.Unlock()

	return ret
}
