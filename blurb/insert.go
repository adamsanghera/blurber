package blurb

import (
	"log"
	"time"

	"github.com/adamsanghera/blurber-protobufs/dist/blurb"
)

// AddNewBlurb consumes metadata to produce a new blurb object,
// and stores that Blurb object in-memory.
func (ll *LocalLedger) AddNewBlurb(creatorID int32, content string, creatorName string) {
	log.Printf("BLURB-LEDGER: Adding a new blurb %d, created by %s", ll.bidCounter, creatorName)

	freshBID := ll.measureBID()

	// Create a new mapping for a user, if one does not exist.
	value, exists := ll.ledger.Load(creatorID)
	if !exists {
		value = NewBox()
		ll.ledger.Store(creatorID, value)
	}
	box, ok := value.(*Box)
	if !ok {
		panic("Something went terribly wrong")
	}

	// Birth the blurb
	creationTime := time.Now()
	box.insert(
		blurb.Blurb{
			Content:     content,
			UnixTime:    creationTime.Unix(),
			Timestamp:   creationTime.Format("Jan 2 – 15:04 EDT"),
			BlurbID:     freshBID,
			CreatorName: creatorName,
		})
}
