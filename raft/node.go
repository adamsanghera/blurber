package raft

import (
	"log"
	"time"
)

type logEntry struct {
	cmd  interface{}
	term int
}

type AppendEntriesRequest struct {
	leaderTerm   int
	leaderID     int
	prevLogIdx   int
	prevLogTerm  int
	entries      []logEntry
	leaderCommit int
}

type VoteRequest struct {
	candidateTerm int
	candidateID   int
	lastLogIdx    int
	lastLogTerm   int
}

// Node implements all of the RPCs needed to be a self-sufficient RAFT node
type Node struct {
	// Persistent state
	log             []logEntry
	currentTerm     int
	votedFor        int
	state           int // 0=follower, 1=candidate, 2=leader
	electionTimeout time.Duration
	heartBeat       chan struct{}
	id              int
	peers           []struct{}

	// Volatile state
	commitIdx      int
	lastAppliedIdx int

	// Leader-specific state
	nextIdx  []int
	matchIdx []int
}

// NewNode initializes and returns a new node
func NewNode(id int) *Node {
	n := &Node{
		// Persistent state
		log:             make([]logEntry, 0),
		currentTerm:     0,
		votedFor:        -1,
		state:           0,
		electionTimeout: time.Millisecond * 300,
		id:              id,

		// Volatile state
		commitIdx:      0,
		lastAppliedIdx: 0,

		// Leader-specific state
		nextIdx:  make([]int, 0),
		matchIdx: make([]int, 0),
	}
	go n.fightTheMan()
	return n
}

// AppendEntries is how a node's log is appended to.
// This is the RPC function called by the RAFT leader.
func (n *Node) AppendEntries(aer *AppendEntriesRequest) (term int, success bool) {
	log.Printf("RAFT: AppendEntries request received")

	if aer.leaderTerm < n.currentTerm {
		return n.currentTerm, false // the leader is delusional
	}
	if aer.prevLogTerm >= len(n.log) || n.log[aer.prevLogIdx].term != aer.prevLogTerm {
		return n.currentTerm, false // this node has fallen behind
	}

	// If an existing entry conflicts with a new one (same index but different terms), delete the existing entry and all that follow it (§5.3)
	entriesIdx := 0
	logIdx := aer.prevLogIdx + 1
	for logIdx < len(n.log) && entriesIdx < len(aer.entries) {
		if n.log[logIdx].term != aer.entries[entriesIdx].term {
			n.log = n.log[:logIdx]
			break
		}

		entriesIdx++
		logIdx++
	}

	// Append any entries not already in log
	entriesIdx = 0
	logIdx = aer.prevLogIdx + 1
	for entriesIdx < len(aer.entries) {
		// extend the whole log
		if logIdx > len(n.log) {
			n.log = append(n.log, aer.entries[entriesIdx:]...)
			logIdx = len(n.log)
			break
		}

		// fill in the log piece by piece
		n.log[logIdx] = aer.entries[entriesIdx]

		entriesIdx++
		logIdx++
	}
	lastIdx := logIdx - 1

	if aer.leaderCommit > n.commitIdx {
		n.commitIdx = min(aer.leaderCommit, lastIdx)
	}

	return n.currentTerm, true
}

func (n *Node) RequestVote(vr *VoteRequest) (term int, voteGranted bool) {
	if vr.candidateTerm < n.currentTerm {
		return n.currentTerm, false
	}

	// Raft determines which of two logs is more up-to-date by comparing the index and term of the last entries in the logs. If the logs have last entries with different terms, then the log with the later term is more up-to-date. If the logs end with the same term, then whichever log is longer is more up-to-date.

	if n.votedFor == -1 || n.votedFor == vr.candidateID {
		if vr.lastLogTerm > n.log[len(n.log)-1].term {
			n.votedFor = vr.candidateID
			return n.currentTerm, true
		}
		if vr.lastLogTerm == n.log[len(n.log)-1].term {
			if vr.lastLogIdx >= len(n.log)-1 {
				n.votedFor = vr.candidateID
				return n.currentTerm, true
			}
		}
	}

	return n.currentTerm, false
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}
