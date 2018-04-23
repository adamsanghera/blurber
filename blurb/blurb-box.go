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
}

// NewBox creates a new Box
func NewBox() *Box {
	return &Box{
		Box:         sync.Map{},
		SortedCache: make([]blurb.Blurb, 0),
	}
}

// Insert
func (bb *Box) insert(b blurb.Blurb) {
	// bb.Box[b.BID] = b
	bb.Box.Store(b.BlurbID, b)

	bb.SortedCache = append([]blurb.Blurb{b}, bb.SortedCache...)

	if len(bb.SortedCache) > 10 {
		bb.SortedCache = bb.SortedCache[:10]
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

func (bb *Box) sortedCache() []blurb.Blurb {
	return bb.SortedCache
}
