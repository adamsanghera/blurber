package blurb

import "log"

// RemoveBlurb deletes a blurb from the ledger.
func (ll *LocalLedger) RemoveBlurb(creatorID int, bid int) {
	log.Printf("BLURB-LEDGER: Deleting blurb %d", bid)
	ll.ledger[creatorID].delete(bid)
}

// RemoveAllBlurbsBy deletes all blurbs created by a given user.
func (ll *LocalLedger) RemoveAllBlurbsBy(creatorID int) {
	log.Printf("BLURB-LEDGER: Removing all blurbs by %d", creatorID)
	delete(ll.ledger, creatorID)
}
