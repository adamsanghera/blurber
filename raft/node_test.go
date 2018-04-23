package raft

import (
	"testing"
)

func TestNode_AppendEntries(t *testing.T) {
	tests := []struct {
		name string
		aer  *AppendEntriesRequest
	}{
		struct {
			name string
			aer  *AppendEntriesRequest
		}{
			"basic test",
			&AppendEntriesRequest{
				leaderTerm:   0,
				leaderID:     1,
				prevLogIdx:   0,
				entries:      []logEntry{logEntry{cmd: "1", term: 0}, logEntry{cmd: "2", term: 0}, logEntry{cmd: "3", term: 0}},
				leaderCommit: 2,
			}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := NewNode()
			n.currentTerm = 3

			term, succ := n.AppendEntries(tt.aer)
			if succ {
				t.Fatalf("Should have been unsuccessful")
			}
			if term != 3 {
				t.Fatalf("Wrong term %d", term)
			}

			t.Fatalf("Log: %v", n.log)

		})
	}
}
