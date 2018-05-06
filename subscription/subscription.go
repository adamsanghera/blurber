package subscription

import (
	"log"
	"sync"
)

// LocalLedger implements the Ledger interface,
// by maintaining an in-memory record of follower-leader relationships.
type LocalLedger struct {
	followLead sync.Map // maps followerID -> map-set of leaderIDs
	leadFollow sync.Map // maps leaderID -> map-set of followerIDs
}

// NewLocalLedger creates a new LocalLedger, instantiating its
// in-memory data structures.
func NewLocalLedger() *LocalLedger {
	log.Printf("SUB-LEDGER: Initializing")
	return &LocalLedger{
		followLead: sync.Map{},
		leadFollow: sync.Map{},
	}
}

// AddSub adds the subscription (followerID -> leaderID) to the ledger
func (ll *LocalLedger) AddSub(followerID int32, leaderID int32) {
	log.Printf("SUB-LEDGER: Add sub (%d->%d)", followerID, leaderID)
	if _, exists := ll.followLead.Load(followerID); !exists {
		ll.followLead.Store(followerID, &sync.Map{})
	}

	if _, exists := ll.leadFollow.Load(leaderID); !exists {
		ll.leadFollow.Store(leaderID, &sync.Map{})
	}

	m, _ := ll.followLead.Load(followerID)
	fl, ok := m.(*sync.Map)
	if !ok {
		panic("SUB-LEDGER: Something bad sneaked into the ledger")
	}

	m, _ = ll.leadFollow.Load(leaderID)
	lf, ok := m.(*sync.Map)
	if !ok {
		panic("SUB-LEDGER: Something bad sneaked into the ledger")
	}

	fl.Store(leaderID, struct{}{})
	lf.Store(followerID, struct{}{})
}

// RemoveSub removes the subscription (followerID -> leaderID) from the ledger
func (ll *LocalLedger) RemoveSub(followerID int32, leaderID int32) {
	log.Printf("SUB-LEDGER: Remove sub (%d->%d)", followerID, leaderID)
	if _, exists := ll.followLead.Load(followerID); exists {
		m, _ := ll.followLead.Load(followerID)
		set, ok := m.(*sync.Map)
		if !ok {
			panic("SUB-LEDGER: Something bad sneaked into the ledger")
		}
		set.Delete(leaderID)
	}

	if _, exists := ll.leadFollow.Load(leaderID); exists {
		m, _ := ll.leadFollow.Load(leaderID)
		lf, ok := m.(*sync.Map)
		if !ok {
			panic("SUB-LEDGER: Something bad sneaked into the ledger")
		}

		lf.Delete(followerID)
	}
}

// RemoveUser removes all subscriptions (X -> uid) and (uid -> X) from the ledger
func (ll *LocalLedger) RemoveUser(uid int32) {
	log.Printf("SUB-LEDGER: Remove %d from all records", uid)

	// Obtain list of leaders
	leaders, err := ll.GetLeaders(uid)
	if err != nil {
		panic("SUB-LEDGER can't get leaders")
	}

	// Range over leaders
	for _, leader := range leaders {
		ll.RemoveSub(uid, leader) // removes from both lists
	}

	// Delete leader and follower lists for uid
	ll.followLead.Delete(uid)
	ll.leadFollow.Delete(uid)
}

// GetLeaders obtains all subscriptions of the form (uid -> X) from the ledger
func (ll *LocalLedger) GetLeaders(uid int32) ([]int32, error) {
	log.Printf("SUB-LEDGER: Getting all leaders for %d", uid)

	ret := make([]int32, 0)

	subs, exists := ll.followLead.Load(uid)
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

func (ll *LocalLedger) GetFollowers(uid int32) ([]int32, error) {
	log.Printf("SUB-LEDGER: Getting all followers for %d", uid)

	ret := make([]int32, 0)

	subs, exists := ll.leadFollow.Load(uid)
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
