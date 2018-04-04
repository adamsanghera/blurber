package blurb

import (
	"log"
	"sync"
)

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
	defer bb.BoxLock.Unlock()

	bb.Box[b.BID] = b

	bb.SortedCache = append([]Blurb{b}, bb.SortedCache...)

	if len(bb.SortedCache) > 10 {
		bb.SortedCache = bb.SortedCache[:9]
	}

	log.Printf("Sorted cache:\n %v", bb.SortedCache)
}

func (bb *BlurbBox) delete(bid int) {
	bb.BoxLock.Lock()
	defer bb.BoxLock.Unlock()

	delete(bb.Box, bid)

	for k, v := range bb.SortedCache {
		if v.BID == bid {
			bb.SortedCache = append(bb.SortedCache[:k], bb.SortedCache[k+1:]...)
			break
		}
	}
}

func (bb *BlurbBox) snapshot() []Blurb {
	bb.BoxLock.Lock()
	defer bb.BoxLock.Unlock()

	blurbs := make([]Blurb, 0)
	for _, v := range bb.Box {
		blurbs = append(blurbs, v)
	}

	return blurbs
}

func (bb *BlurbBox) sortedCache() []Blurb {
	bb.BoxLock.Lock()
	defer bb.BoxLock.Unlock()

	return bb.SortedCache
}
