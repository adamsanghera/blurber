package blurb

import "log"

// GetBlurbsCreatedBy returns all blurbs created by a given user.
func (ll *LocalLedger) GetBlurbsCreatedBy(creatorID int) []Blurb {
	log.Printf("BLURB-LEDGER: Retrieving all blurbs by %d", creatorID)

	return ll.ledger[creatorID].snapshot()
}
