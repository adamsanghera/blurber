package blurb

import "log"

// RemoveBlurb deletes a blurb from the ledger.
func (ll *LocalLedger) RemoveBlurb(creatorID int32, bid int32) {
	log.Printf("BLURB-LEDGER: Deleting blurb %d", bid)
	val, exists := ll.ledger.Load(creatorID)
	if !exists {
		return
	}

	box := val.(*Box)
	box.delete(bid)
}

// RemoveAllBlurbsBy deletes all blurbs created by a given user.
func (ll *LocalLedger) RemoveAllBlurbsBy(creatorID int32) {
	log.Printf("BLURB-LEDGER: Removing all blurbs by %d", creatorID)
	ll.ledger.Delete(creatorID)
	ll.feeds.cache.Delete(creatorID)
}
