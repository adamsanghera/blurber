package blurb

import "sync"

type BlurbBox struct {
	Box         map[int]Blurb
	SortedCache []Blurb
	BoxLock     sync.Mutex
}

func NewBlurbBox() *BlurbBox {
	return &BlurbBox{
		Box:         make(map[int]Blurb),
		BoxLock:     sync.Mutex{},
		SortedCache: make([]Blurb, 0),
	}
}

func (bb *BlurbBox) insert(b Blurb) {
	bb.BoxLock.Lock()

	bb.Box[b.BID] = b

	bb.BoxLock.Unlock()
}

func (bb *BlurbBox) delete(bid int) {
	bb.BoxLock.Lock()

	delete(bb.Box, bid)

	bb.BoxLock.Unlock()
}

func (bb *BlurbBox) snapshot() []Blurb {
	bb.BoxLock.Lock()

	blurbs := make([]Blurb, 0)
	for _, v := range bb.Box {
		blurbs = append(blurbs, v)
	}

	bb.BoxLock.Unlock()

	return blurbs
}
