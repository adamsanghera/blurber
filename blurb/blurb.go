package blurb
import "log"
type Blurb struct {
	Content     string
	Timestamp   string
	BID         string // immutable
	CreatorName string
}

type BlurbLedger interface {
	AddBlurb(creator string, b Blurb)
	RemoveBlurb(creator string, b Blurb)
}

type LocalBlurbLedger struct {
	// Map from UID -> map [BID] -> Blurb
	ledger map[string]map[string]Blurb
}

func NewLocalLedger() *LocalBlurbLedger {
	return &LocalBlurbLedger{
		ledger:  make(map[string]map[string]Blurb),
	}
}

func (lbl LocalBlurbLedger) AddBlurb(creator string, b Blurb) {
	log.Printf("Adding %v %v...", creator, b)
	if _, exists := lbl.ledger[creator][b.BID]; !exists {
		lbl.ledger[creator] = make(map[string]Blurb)
		log.Printf("%v", lbl.ledger)
	}
	lbl.ledger[creator][b.BID] = b
	log.Printf("%v", lbl.ledger)
}

func (lbl LocalBlurbLedger) RemoveBlurb(creator string, BID string) {
	delete(lbl.ledger[creator], BID)
}

func (lbl LocalBlurbLedger) GetUsrBlurb(creator string) []Blurb {
	bs := make([]Blurb, 0)
	for _, v := range lbl.ledger[creator] {
		bs = append(bs, v)
	}
	return bs
}