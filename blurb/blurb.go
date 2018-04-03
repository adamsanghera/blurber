package blurb

import (
	"log"
	"time"
)

type Blurb struct {
	content     string
	timestamp   string
	bid         int // immutable
	creatorName string
}

type BlurbLedger interface {
	AddBlurb(creator string, b Blurb)
	RemoveBlurb(creator string, b Blurb)
}

type LocalBlurbLedger struct {
	// Map from UID -> map [BID] -> Blurb
	bidCounter int
	ledger     map[int]map[int]Blurb
}

func NewLocalLedger() *LocalBlurbLedger {
	return &LocalBlurbLedger{
		ledger:     make(map[int]map[int]Blurb),
		bidCounter: 0,
	}
}

func (lbl LocalBlurbLedger) AddNewBlurb(creatorID int, content string, creatorName string) {
	if _, exists := lbl.ledger[creatorID]; !exists {
		lbl.ledger[creatorID] = make(map[int]Blurb)
	}
	lbl.ledger[creatorID][lbl.bidCounter] = Blurb{
		content:     content,
		timestamp:   time.Now().Format("Jan 2 – 15:04 EDT"),
		bid:         lbl.bidCounter,
		creatorName: creatorName,
	}
	lbl.bidCounter++
	log.Printf("Updated lbl: %v", lbl.ledger)
}

func (lbl LocalBlurbLedger) RemoveBlurb(creatorID int, bid int) {
	delete(lbl.ledger[creatorID], bid)
}

func (lbl LocalBlurbLedger) GetUsrBlurb(creatorID int) []Blurb {
	bs := make([]Blurb, 0)
	for _, v := range lbl.ledger[creatorID] {
		bs = append(bs, v)
	}
	return bs
}
