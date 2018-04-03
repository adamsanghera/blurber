package subscription

import (
	"log"
)

type SubscriptionLedger interface {
	AddSub(followerID int, leadearID int)
	RemoveSub(followerID int, leadearID int)
}

type LocalSubscriptionLedger struct {
	// Map from UID-> map[UID] -> nonce
	ledger map[int]map[int]bool
}

func NewLocalLedger() *LocalSubscriptionLedger {
	return &LocalSubscriptionLedger{
		ledger: make(map[int]map[int]bool),
	}
}

func (lsl *LocalSubscriptionLedger) AddSub(followerID int, leadearID int) {
	if _, exists := lsl.ledger[followerID]; !exists {
		lsl.ledger[followerID] = make(map[int]bool)
	}
	//
	lsl.ledger[followerID][leadearID] = true
	log.Printf("Updated lsl: %v", lsl.ledger)
}

func (lsl *LocalSubscriptionLedger) RemoveSub(followerID int, leadearID int) {
	if _, exists := lsl.ledger[followerID]; !exists {
		return
	}
	delete(lsl.ledger[followerID], leadearID)
}

func (lsl *LocalSubscriptionLedger) RemoveUsr(uid int) {
	log.Printf("SUBSCRIPTION: Removing %s and all their followers", uid)

	// Unsubscribe this person from every follower in the ledger
	for followID := range lsl.ledger {
		lsl.RemoveSub(followID, uid)
	}

	// Remove his own subscription list
	delete(lsl.ledger, uid)
}
