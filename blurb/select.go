package blurb

import (
	"log"
)

// GetBlurbsCreatedBy returns all blurbs created by a given user.
func (ll *LocalLedger) GetBlurbsCreatedBy(creatorID int32) []Blurb {
	log.Printf("BLURB-LEDGER: Retrieving all blurbs by %d", creatorID)

	if _, exists := ll.ledger[creatorID]; !exists {
		return make([]Blurb, 0)
	}

	return ll.ledger[creatorID].snapshot()
}

func (ll *LocalLedger) GetRecentBlurbsBy(creatorID int32) []Blurb {
	log.Printf("BLURB-LEDGER: Retrieving blurb cache for %d", creatorID)

	if _, exists := ll.ledger[creatorID]; !exists {
		return make([]Blurb, 0)
	}

	return ll.ledger[creatorID].sortedCache()
}
