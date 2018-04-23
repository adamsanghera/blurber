package raft

import "context"

func (n *Node) fightTheMan() {
	ctx, _ := context.WithTimeout(nil, n.electionTimeout)

	select {
	case <-ctx.Done():
		// begin vote protocol
		n.ClaimTheThrone()
		break
	case <-n.heartBeat:
		ctx, _ = context.WithTimeout(nil, n.electionTimeout)
	}
}

func (n *Node) ClaimTheThrone() {
	n.currentTerm++
	n.state = 1
	n.votedFor = n.id

	for range n.peers {
		n.RequestVote(&VoteRequest{
			candidateTerm: n.currentTerm,
			candidateID:   n.id,
			lastLogIdx:    n.lastAppliedIdx,
			lastLogTerm:   n.log[len(n.log)-1].term, // This might be wrong
		})
	}
}
