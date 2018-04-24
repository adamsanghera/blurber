package blurb

import (
	"log"

	"github.com/adamsanghera/blurber-protobufs/dist/blurb"
)

// GetBlurbsCreatedBy returns all blurbs created by a given user.
func (ll *LocalLedger) GetBlurbsCreatedBy(creatorID int32) []blurb.Blurb {
	log.Printf("BLURB-LEDGER: Retrieving all blurbs by %d", creatorID)

	value, exists := ll.ledger.Load(creatorID)
	if !exists {
		return make([]blurb.Blurb, 0)
	}

	box, ok := value.(*Box)
	if !ok {
		panic("Something bad was inserted into the ledger")
	}

	return box.snapshot()
}

// GetRecentBlurbsBy returns the N most recent blurbs made by the creatorID
func (ll *LocalLedger) GetRecentBlurbsBy(creatorID int32) []blurb.Blurb {
	log.Printf("BLURB-LEDGER: Retrieving blurb cache for %d", creatorID)

	value, exists := ll.ledger.Load(creatorID)
	if !exists {
		return make([]blurb.Blurb, 0)
	}

	box, ok := value.(*Box)
	if !ok {
		panic("Something bad was inserted into the ledger")
	}

	return box.getCache()
}
