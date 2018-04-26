package blurb

import (
	"log"
	"sync"

	"github.com/adamsanghera/blurber-protobufs/dist/blurb"
)

// Box is where blurbs are stored.
// There is one Box for every user.
type Box struct {
	Box         sync.Map
	SortedCache []blurb.Blurb

	cacheSize int
}

// NewBox creates a new Box
func NewBox(cacheSize int) *Box {
	return &Box{
		Box:         sync.Map{},
		SortedCache: make([]blurb.Blurb, 0),

		cacheSize: cacheSize,
	}
}

// Insert
func (bb *Box) insert(b blurb.Blurb) {
	// bb.Box[b.BID] = b
	bb.Box.Store(b.BlurbID, b)

	bb.SortedCache = append([]blurb.Blurb{b}, bb.SortedCache...)

	if len(bb.SortedCache) > bb.cacheSize {
		bb.SortedCache = bb.SortedCache[:bb.cacheSize]
	}

	log.Printf("Sorted cache:\n %v", bb.SortedCache)
}

func (bb *Box) delete(bid int32) {
	bb.Box.Delete(bid)

	for k, v := range bb.SortedCache {
		if v.BlurbID == bid {
			bb.SortedCache = append(bb.SortedCache[:k], bb.SortedCache[k+1:]...)
			break
		}
	}
}

func (bb *Box) snapshot() []blurb.Blurb {
	blurbs := make([]blurb.Blurb, 0)

	bb.Box.Range(func(k interface{}, v interface{}) bool {
		_, ok := k.(int32)
		if !ok {
			return false
		}
		val, ok := v.(blurb.Blurb)
		if !ok {
			return false
		}

		blurbs = append(blurbs, val)
		return true
	})

	return blurbs
}

func (bb *Box) getCache() []blurb.Blurb {
	return bb.SortedCache
}
