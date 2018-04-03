package blurb

import (
	"log"
	"time"
)

type Blurb struct {
	Content     string
	Timestamp   string
	Time        time.Time
	BID         int // immutable
	CreatorName string
}

// Ledger is the interface for a service that maintains
// a record of blurbs (social media posts)
type Ledger interface {
	AddBlurb(creator string, b Blurb)
	RemoveBlurb(creator string, b Blurb)
}

// LocalLedger is an implementation of Ledger, making use
// of in-memory data structures to record blurbs
type LocalLedger struct {
	// Map from UID -> map [BID] -> Blurb
	bidCounter int
	ledger     map[int]map[int]Blurb
}

// NewLocalLedger returns a new and initialized LocalLedger
func NewLocalLedger() *LocalLedger {
	log.Printf("BLURB-LEDGER: Initializing")
	return &LocalLedger{
		ledger:     make(map[int]map[int]Blurb),
		bidCounter: 0,
	}
}

// AddNewBlurb consumes metadata to produce a new blurb object,
// and stores that Blurb object in-memory.
func (lbl *LocalLedger) AddNewBlurb(creatorID int, content string, creatorName string) {
	log.Printf("BLURB-LEDGER: Adding a new blurb %d, created by %s", lbl.bidCounter, creatorName)

	// Create a new mapping for a user, if one does not exist.
	if _, exists := lbl.ledger[creatorID]; !exists {
		lbl.ledger[creatorID] = make(map[int]Blurb)
	}

	// Birth the blurb
	creationTime := time.Now()
	lbl.ledger[creatorID][lbl.bidCounter] = Blurb{
		Content:     content,
		Time:        creationTime,
		Timestamp:   creationTime.Format("Jan 2 – 15:04 EDT"),
		BID:         lbl.bidCounter,
		CreatorName: creatorName,
	}

	// Increment the BID counter
	lbl.bidCounter++
}

// RemoveBlurb deletes a blurb from the ledger.
func (lbl *LocalLedger) RemoveBlurb(creatorID int, bid int) {
	log.Printf("BLURB-LEDGER: Deleting blurb %d", bid)
	delete(lbl.ledger[creatorID], bid)
}

// GetBlurbsCreatedBy returns all blurbs created by a given user.
func (lbl *LocalLedger) GetBlurbsCreatedBy(creatorID int) []Blurb {
	log.Printf("BLURB-LEDGER: Retrieving all blurbs by %d", creatorID)
	bs := make([]Blurb, 0)
	for _, v := range lbl.ledger[creatorID] {
		bs = append(bs, v)
	}
	return bs
}

// RemoveAllBlurbsBy deletes all blurbs created by a given user.
func (lbl *LocalLedger) RemoveAllBlurbsBy(creatorID int) {
	log.Printf("BLURB-LEDGER: Removing all blurbs by %d", creatorID)
	for k := range lbl.ledger[creatorID] {
		lbl.RemoveBlurb(creatorID, k)
	}
}
