package blurb

import (
	"log"
	"time"
)

// AddNewBlurb consumes metadata to produce a new blurb object,
// and stores that Blurb object in-memory.
func (ll *LocalLedger) AddNewBlurb(creatorID int, content string, creatorName string) {
	log.Printf("BLURB-LEDGER: Adding a new blurb %d, created by %s", ll.bidCounter, creatorName)

	freshBID := ll.measureBID()

	// Create a new mapping for a user, if one does not exist.
	if _, exists := ll.ledger[creatorID]; !exists {
		ll.ledger[creatorID] = NewBlurbBox()
	}

	// Birth the blurb
	creationTime := time.Now()
	ll.ledger[creatorID].insert(
		Blurb{
			Content:     content,
			Time:        creationTime,
			Timestamp:   creationTime.Format("Jan 2 – 15:04 EDT"),
			BID:         freshBID,
			CreatorName: creatorName,
		})
}
