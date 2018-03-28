package main

import "time"

type Blurb struct {
	Content    string
	Timestamp  time.Time
	BID        string // immutable
	CreatorUID string
}

type BlurbLedger interface {
	AddBlurb(creator string, b Blurb)
	RemoveBlurb(creator string, b Blurb)
}

type LocalBlurbLedger struct {
	// Map from UID -> map [BID] -> Blurb
	ledger map[string]map[string]Blurb
}

func (lbl LocalBlurbLedger) AddBlurb(creator string, b Blurb) {
	if _, exists := lbl.ledger[creator][b.BID]; !exists {
		lbl.ledger[creator] = make(map[string]Blurb)
	}
	lbl.ledger[creator][b.BID] = b
}

func (lbl LocalBlurbLedger) RemoveBlurb(creator string, BID string) {
	delete(lbl.ledger[creator], BID)
}
