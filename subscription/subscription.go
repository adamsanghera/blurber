package subscription

import (
	"log"
	"sync"
)

// Ledger is the interface for a service that maintains a
// record of follower-leader relationships in a social graph.
type Ledger interface {
	AddSub(followerID int32, leadearID int32)
	RemoveSub(followerID int32, leadearID int32)
}

// LocalLedger implements the Ledger interface,
// by maintaining an in-memory record of follower-leader relationships.
type LocalLedger struct {
	ledger sync.Map // maps followerID -> map-set of leaderIDs
}

// NewLocalLedger creates a new LocalLedger, instantiating its
// in-memory data structures.
func NewLocalLedger() *LocalLedger {
	log.Printf("SUB-LEDGER: Initializing")
	return &LocalLedger{
		ledger: sync.Map{},
	}
}

// AddSub adds the subscription (followerID -> leaderID) to the ledger
func (ll *LocalLedger) AddSub(followerID int32, leaderID int32) {
	log.Printf("SUB-LEDGER: Add sub (%d->%d)", followerID, leaderID)
	if _, exists := ll.ledger.Load(followerID); !exists {
		ll.ledger.Store(followerID, &sync.Map{})
	}

	m, _ := ll.ledger.Load(followerID)
	set, ok := m.(*sync.Map)
	if !ok {
		panic("SUB-LEDGER: Something bad sneaked into the ledger")
	}

	set.Store(leaderID, struct{}{})
}

// RemoveSub removes the subscription (followerID -> leaderID) from the ledger
func (ll *LocalLedger) RemoveSub(followerID int32, leaderID int32) {
	log.Printf("SUB-LEDGER: Remove sub (%d->%d)", followerID, leaderID)
	if _, exists := ll.ledger.Load(followerID); !exists {
		return
	}

	m, _ := ll.ledger.Load(followerID)
	set, ok := m.(*sync.Map)
	if !ok {
		panic("SUB-LEDGER: Something bad sneaked into the ledger")
	}

	set.Delete(leaderID)
}

// RemoveUser removes all subscriptions (X -> uid) and (uid -> X) from the ledger
func (ll *LocalLedger) RemoveUser(uid int32) {
	log.Printf("SUB-LEDGER: Remove %d from all records", uid)

	// Unsubscribe this person from every follower in the ledger
	ll.ledger.Range(func(key interface{}, _ interface{}) bool {
		keyInt, ok := key.(int32)
		if !ok {
			panic("SUB-LEDGER: Bad key in ledger")
		}
		ll.RemoveSub(keyInt, uid)
		return true
	})

	// Remove his own subscription list
	ll.ledger.Delete(uid)
}

// GetLeaders obtains all subscriptions of the form (uid -> X) from the ledger
func (ll *LocalLedger) GetLeaders(uid int32) ([]int32, error) {
	log.Printf("SUB-LEDGER: Getting all leaders for %d", uid)

	ret := make([]int32, 0)

	subs, exists := ll.ledger.Load(uid)
	if !exists {
		return ret, nil
	}

	subMap, ok := subs.(*sync.Map)
	if !ok {
		panic("SUB-LEDGER: Bad value in ledger")
	}

	subMap.Range(func(key interface{}, _ interface{}) bool {
		id, ok := key.(int32)
		if !ok {
			panic("SUB-LEDGER: Bad value in ledger")
		}
		ret = append(ret, id)
		return true
	})

	return ret, nil
}
