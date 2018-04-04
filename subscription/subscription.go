package subscription

import (
	"log"
)

// Ledger is the interface for a service that maintains a
// record of follower-leader relationships in a social graph.
type Ledger interface {
	AddSub(followerID int, leadearID int)
	RemoveSub(followerID int, leadearID int)
}

// LocalLedger implements the Ledger interface,
// by maintaining an in-memory record of follower-leader relationships.
type LocalLedger struct {
	// Map from UID-> map[UID] -> nonce
	ledger map[int]map[int]bool
}

// NewLocalLedger creates a new LocalLedger, instantiating its
// in-memory data structures.
func NewLocalLedger() *LocalLedger {
	log.Printf("SUB-LEDGER: Initializing")
	return &LocalLedger{
		ledger: make(map[int]map[int]bool),
	}
}

// AddSub adds the subscription (followerID -> leaderID) to the ledger
func (ll *LocalLedger) AddSub(followerID int, leaderID int) {
	log.Printf("SUB-LEDGER: Add sub (%d->%d)", followerID, leaderID)
	if _, exists := ll.ledger[followerID]; !exists {
		ll.ledger[followerID] = make(map[int]bool)
	}

	ll.ledger[followerID][leaderID] = true
}

// RemoveSub removes the subscription (followerID -> leaderID) from the ledger
func (ll *LocalLedger) RemoveSub(followerID int, leaderID int) {
	log.Printf("SUB-LEDGER: Remove sub (%d->%d)", followerID, leaderID)
	if _, exists := ll.ledger[followerID]; !exists {
		return
	}
	delete(ll.ledger[followerID], leaderID)
}

// RemoveUser removes all subscriptions (X -> uid) and (uid -> X) from the ledger
func (ll *LocalLedger) RemoveUser(uid int) {
	log.Printf("SUB-LEDGER: Remove %d from all records", uid)

	// Unsubscribe this person from every follower in the ledger
	for followID := range ll.ledger {
		ll.RemoveSub(followID, uid)
	}

	// Remove his own subscription list
	delete(ll.ledger, uid)
}

// GetLeaders obtains all subscriptions of the form (uid -> X) from the ledger
func (ll *LocalLedger) GetLeaders(uid int) (map[int]bool, error) {
	return ll.ledger[uid], nil
}
