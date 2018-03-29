package subscription

type SubscriptionLedger interface {
	AddSub(follower string, leader string)
	RemoveSub(follower string, leader string)
}

type LocalSubscriptionLedger struct {
	// Map from UID-> map[UID] -> nonce
	ledger map[string]map[string]bool
}

func (lsl LocalSubscriptionLedger) AddSub(follower string, leader string) {
	if _, exists := lsl.ledger[follower]; !exists {
		lsl.ledger[follower] = make(map[string]bool)
	}
	lsl.ledger[follower][leader] = true
}

func (lsl LocalSubscriptionLedger) RemoveSub(follower string, leader string) {
	if _, exists := lsl.ledger[follower]; !exists {
		return
	}
	delete(lsl.ledger[follower], leader)
}
