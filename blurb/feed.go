package blurb

import (
	"log"
	"os"
	"sort"
	"time"

	"github.com/adamsanghera/blurber-protobufs/dist/blurb"
)

func (ll *LocalLedger) GenerateFeed(readerID int32, sources []int32) []blurb.Blurb {
	log.Printf("BLURB-LEDGER: Generating a feed for %d", readerID)

	// Hit cache, if it exists and is fresh
	value, exists := ll.feedCache.Load(readerID)
	if exists {
		cache, ok := value.(FeedCache)
		if !ok {
			panic("Something very wrong snuck into feed cache map")
		}

		if time.Since(cache.timestamp) < time.Second*30 {
			log.Printf("BLURB-LEDGER: Hitting the cache for a feed for %d", readerID)
			return cache.sortedFeed
		}
	}

	log.Printf("BLURB-LEDGER: Generating a new feed for %d", readerID)

	// Get recent blurbs from all sources
	ret := make([]blurb.Blurb, 0)
	for _, leaderID := range sources {
		ret = append(ret, ll.GetRecentBlurbsBy(leaderID)...)
	}

	// Sort the blurbs
	sort.Slice(ret, func(i, j int) bool {
		iTime := time.Unix(ret[i].UnixTime, 0)
		jTime := time.Unix(ret[j].UnixTime, 0)
		return iTime.After(jTime)
	})

	if len(ret) > 100 {
		ret = ret[:100]
	}

	// Store the result in the cache
	ll.feedCache.Store(readerID, FeedCache{
		sortedFeed: ret,
		timestamp:  time.Now(),
	})

	if os.Getenv("DEBUG") == "1" {
		log.Printf("BLURB-LEDGER: Feed generated with %d elements", len(ret))
	}

	return ret
}

// InvalidateCache wipes the cache for the given user.
func (ll *LocalLedger) InvalidateCache(readerID int32) {
	log.Printf("BLURB-LEDGER: Invalidating cache for %d", readerID)
	ll.feedCache.Delete(readerID)
}
