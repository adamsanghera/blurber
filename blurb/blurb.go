package blurb

import (
	"log"
	"time"
)

type blurb struct {
	content     string
	timestamp   string
	bid         int // immutable
	creatorName string
}

type BlurbLedger interface {
	AddBlurb(creator string, b blurb)
	RemoveBlurb(creator string, b blurb)
}

type LocalBlurbLedger struct {
	// Map from UID -> map [BID] -> Blurb
	bidCounter int
	ledger     map[string]map[int]blurb
}

func NewLocalLedger() *LocalBlurbLedger {
	return &LocalBlurbLedger{
		ledger:     make(map[string]map[int]blurb),
		bidCounter: 0,
	}
}

func (lbl LocalBlurbLedger) AddNewBlurb(creator string, content string, creatorName string) {
	if _, exists := lbl.ledger[creator]; !exists {
		lbl.ledger[creator] = make(map[int]blurb)
	}
	lbl.ledger[creator][lbl.bidCounter] = blurb{
		content:     content,
		timestamp:   time.Now().Format("Jan 2 – 15:04 EDT"),
		bid:         lbl.bidCounter,
		creatorName: creatorName,
	}
	lbl.bidCounter++
	log.Printf("Updated lbl: %v", lbl.ledger)
}

func (lbl LocalBlurbLedger) RemoveBlurb(creator string, bid int) {
	delete(lbl.ledger[creator], bid)
}

func (lbl LocalBlurbLedger) GetUsrBlurb(creator string) []blurb {
	bs := make([]blurb, 0)
	for _, v := range lbl.ledger[creator] {
		bs = append(bs, v)
	}
	return bs
}
