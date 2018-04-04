package blurb

import (
	"log"
	"os"
	"sort"
	"time"
)

func (ll *LocalLedger) GenerateFeed(readerID int, sources []int) []Blurb {
	log.Printf("BLURB-LEDGER: Generating a feed for %d", readerID)

	// Hit cache, if it exists and is fresh
	if cache, exists := ll.feedCache[readerID]; exists && time.Since(cache.timestamp) < time.Second*30 {
		log.Printf("BLURB-LEDGER: Hitting the cache for a feed for %d", readerID)
		return cache.sortedFeed
	}

	log.Printf("BLURB-LEDGER: Generating a new feed for %d", readerID)

	// Get recent blurbs from all sources
	ret := make([]Blurb, 0)
	for _, leaderID := range sources {
		ret = append(ret, ll.GetRecentBlurbsBy(leaderID)...)
	}

	// Sort the blurbs
	sort.Slice(ret, func(i, j int) bool {
		return ret[i].Time.After(ret[j].Time)
	})

	if len(ret) > 100 {
		ret = ret[:100]
	}

	// Store the result in the cache
	ll.feedCache[readerID] = FeedCache{
		sortedFeed: ret,
		timestamp:  time.Now(),
	}

	if os.Getenv("DEBUG") == "1" {
		log.Printf("BLURB-LEDGER: Feed generated with %d elements", len(ret))
	}

	return ret
}

// InvalidateCache wipes the cache for the given user.
func (ll *LocalLedger) InvalidateCache(readerID int) {
	log.Printf("BLURB-LEDGER: Invalidating cache for %d", readerID)
	delete(ll.feedCache, readerID)
}
